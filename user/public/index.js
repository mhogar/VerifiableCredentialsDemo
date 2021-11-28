const app = {
    template: `
        <div class="container">
            <div v-if="isLoading" class="d-flex justify-content-center">
                <div class="spinner-border" role="status"></div>
            </div>
            <div v-else>
                <div v-if="alert" :class="'text-center alert alert-' + alert.type" role="alert">
                    {{alert.text}}
                </div>
                <div v-if="verifyPrompt" id="verifyPromptCard" class="card text-center">
                    <div class="card-body">
                        <h5 class="card-title">{{verifyPrompt.name}}</h5>
                        <h6 class="card-subtitle mb-2 text-muted">{{verifyPrompt.domain}}</h6>
                        <h6>Trusted By Issuer: {{verifyPrompt.trusted_by_issuer}}</h6>
                        <p class="card-text">{{verifyPrompt.purpose}}</p>
                        <div class="btn-group" role="group">
                            <button type="button" class="btn btn-success" @click="acceptVerifyPromptClick()">Accept</button>
                            <button type="button" class="btn btn-danger" @click="denyVerifyPromptClick">Deny</button>
                        </div>
                    </div>
                </div>
                <form v-else>
                    <!--div class="row">
                        <input id="url-input" class="form-control" v-model="url">   
                    </div-->
                    <div class="row">
                        <button type="submit" class="btn btn-primary" @click.prevent="submitPresReq">Submit</button>
                    </div>
                </form>
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
                    this.setAlert("danger", "An error occurred. Try again later.")
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
                    this.setAlert("success", "Verified!")
                })
                .catch(() => {
                    this.setAlert("danger", "An error occurred. Try again later.")
                })
                .then(() => {
                    this.isLoading = false
                    this.verifyPrompt = false
                })
        }
    }
};

Vue.createApp(app).mount('#app')
