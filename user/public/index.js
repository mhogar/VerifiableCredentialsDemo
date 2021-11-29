const VerifyPage = {
    template: `
        <div>
            <div v-if="verifyPrompt" class="ui raised card">
                <div class="content">
                    <div class="header">{{verifyPrompt.name}}</div>
                    <div class="meta">
                        <span>{{verifyPrompt.domain}}</span>
                    </div>
                    <h4 class="ui sub header">
                        Trusted By Issuer:
                        <i :class="trustedByVerifierIcon"></i>
                    </h4>
                </div>
                <div class="content">
                    <div class="description">
                        <p>{{verifyPrompt.purpose}}</p>
                    </div>
                </div>
                <div class="extra content">
                    <div class="ui buttons">
                        <button type="button" class="ui positive button" @click="acceptVerifyPromptClick()">Accept</button>
                        <div class="or"></div>
                        <button type="button" class="ui negative button" @click="denyVerifyPromptClick">Deny</button>
                    </div>
                </div>
            </div>
            <div v-else>
                <h2 class="ui header">Send a Verify Request</h2>
                <div class="ui action input">
                    <input type="text" v-model="url">
                    <button class="ui button" @click.prevent="submitPresReq">Submit</button>
                </div>
            </div>
        </div>
    `,
    data() {
        return {
            verifyPrompt: null,
            url: "http://localhost:8082/verify"
        }
    },
    computed: {
        trustedByVerifierIcon() {
            return this.verifyPrompt.trusted_by_issuer ? 'check circle green icon' : 'close red icon'
        }
    },
    methods: {
        submitPresReq() {
            this.$emit('clearAlert')
            this.$emit('loading', true)

            axios.get('/verify', {
                params: {
                    url: this.url
                },
                validateStatus: (status) => status < 500
            })
            .then((res) => {
                if (res.data.error) {
                    this.$emit('setAlert', "negative", res.data.error)
                    return
                }

                this.verifyPrompt = res.data
            })
            .catch((err) => {
                console.log(err)
                this.$emit('setAlert', "negative", "An internal error occurred. Try again later.")
            })
            .then(() => {
                this.$emit('loading', false)
            })
        },
        denyVerifyPromptClick() {
            this.$parent.setAlert("warning", "Verify request denied.")
            this.verifyPrompt = null
            this.url = ""
        },
        acceptVerifyPromptClick() {
            this.$emit('loading', true)
            axios.post('/verify', {
                url: this.url
            })
            .then((res) => {
                if (res.data.error) {
                    this.$emit('setAlert', "negative", res.data.error)
                    return
                }

                this.$emit('setAlert', "positive", "Verified!")
            })
            .catch((err) => {
                console.log(err)
                this.$emit('setAlert', "negative", "An internal error occurred. Try again later.")
            })
            .then(() => {
                this.$emit('loading', false)
                this.verifyPrompt = null
                this.url = ""
            })
        }
    }
}

const IssuePage = {
    template: `
        <div>
            <div v-if="issuePrompt" class="ui raised card">
                <div class="content">
                    <div class="header">{{issuePrompt.name}}</div>
                    <div class="meta">
                        <span>{{issuePrompt.domain}}</span>
                    </div>
                </div>
                <div class="content">
                    <div class="description">
                        <p>{{issuePrompt.purpose}}</p>
                    </div>
                </div>
                <div class="extra content">
                    <div class="ui buttons">
                        <button type="button" class="ui positive button" @click="proceedIssuePromptClick()">Proceed</button>
                        <div class="or"></div>
                        <button type="button" class="ui negative button" @click="abortIssuePromptClick">Abort</button>
                    </div>
                </div>
            </div>
            <div v-else-if="fieldsForm">
                <div v-for="(value,key) in fieldsForm">
                    <div class="ui fluid labeled input formInput">
                        <div class="ui label">{{key}}</div>
                        <input type="text" v-model="fieldsForm[key]">
                    </div>
                </div>
                <button class="ui button" @click="submitIssueButtonClick">Submit</button>
            </div>
            <div v-else>
                <h2 class="ui header">Send a Issue Request</h2>
                <div class="ui action input">
                    <input type="text" v-model="url">
                    <button class="ui button" @click="submitIssReq">Submit</button>
                </div>
            </div>
        </div>
    `,
    data() {
        return {
            issuePrompt: null,
            fieldsForm: null,
            url: "http://localhost:8082/issue"
        }
    },
    methods: {
        submitIssReq() {
            this.$emit('clearAlert')
            this.$emit('loading', true)

            axios.get('/issue', {
                params: {
                    url: this.url
                },
                validateStatus: (status) => status < 500
            })
            .then((res) => {
                if (res.data.error) {
                    this.$emit('setAlert', "negative", res.data.error)
                    return
                }

                this.issuePrompt = res.data
            })
            .catch((err) => {
                console.log(err)
                this.$emit('setAlert', "negative", "An internal error occurred. Try again later.")
            })
            .then(() => {
                this.$emit('loading', false)
            })
        },
        abortIssuePromptClick() {
            this.$parent.setAlert("warning", "Issue request aborted.")
            this.issuePrompt = null
            this.url = ""
        },
        proceedIssuePromptClick() {
            this.fieldsForm = this.issuePrompt.fields
            this.issuePrompt = null
        },
        submitIssueButtonClick() {
            this.$emit('loading', true)
            axios.post('/issue', {
                url: this.url,
                fields: this.fieldsForm
            })
            .then((res) => {
                if (res.data.error) {
                    this.$emit('setAlert', "negative", res.data.error)
                    return
                }

                this.$emit('setAlert', "positive", "Verifiable Credentials Created!")
            })
            .catch((err) => {
                console.log(err)
                this.$emit('setAlert', "negative", "An internal error occurred. Try again later.")
            })
            .then(() => {
                this.$emit('loading', false)
                this.fieldsForm = null
                this.url = ""
            })
        }
    }
}

const app = {
    template: `
        <div id="navbar" class="ui fixed borderless huge inverted menu">
            <div class="header item"><b>VCD</b></div>
            <a class="item">Issue</a>
            <a class="item">Verify</a>
        </div>
        <div class="ui basic segment">
            <div v-if="isLoading" class="ui active loader"></div>
            <div v-else class="ui centered container grid">
                <div v-if="alert" id="alert-box" :class="'ui message ' + this.alert.type">
                    <p>
                        {{alert.text}} 
                        <i class="close icon" @click="clearAlert()"></i>
                    </p>
                </div>
                <div id="page-content">
                    <IssuePage @loading="changeLoading" @setAlert="setAlert" @clearAlert="clearAlert" />
                </div>
            </div>
        </div>
    `,
    data() {
        return {
            isLoading: false,
            alert: null,
        }
    },
    components: {
        VerifyPage,
        IssuePage
    },
    methods: {
        changeLoading(loading) {
            this.loading = loading
        },
        setAlert(type, text) {
            this.alert = {
                type: type,
                text: text
            }
        },
        clearAlert() {
            this.alert = null
        }
    }
}

Vue.createApp(app).mount('#app')
