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
                    <div v-if="verifyPrompt" id="verifyPromptCard" class="ui raised card">
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
            </div>
        </div>
    `,
    data() {
        return {
            isLoading: false,
            alert: null,
            verifyPrompt: null,
            url: ""
        }
    },
    computed: {
        trustedByVerifierIcon() {
            return this.verifyPrompt.trusted_by_issuer ? 'check circle green icon' : 'close red icon'
        }
    },
    methods: {
        setAlert(type, text) {
            this.alert = {
                type: type,
                text: text
            }
        },
        clearAlert() {
            this.alert = null
        },
        submitPresReq() {
            this.clearAlert()
            this.isLoading = true

            axios.get('/verify')
                .then((res) => {
                    this.verifyPrompt = res.data
                })
                .catch(() => {
                    this.setAlert("negative", "An error occurred. Try again later.")
                })
                .then(() => {
                    this.isLoading = false
                })
        },
        denyVerifyPromptClick() {
            this.setAlert("warning", "Verify request denied.")
            this.verifyPrompt = false
        },
        acceptVerifyPromptClick() {
            this.isLoading = true
            axios.post('/verify', {})
                .then(() => {
                    this.setAlert("positive", "Verified!")
                })
                .catch(() => {
                    this.setAlert("negative", "An error occurred. Try again later.")
                })
                .then(() => {
                    this.isLoading = false
                    this.verifyPrompt = false
                })
        }
    }
};

Vue.createApp(app).mount('#app')
