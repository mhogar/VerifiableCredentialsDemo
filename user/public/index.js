const app = {
    template: `
        <div class="container d-flex justify-content-center">
            <div v-if="isLoading" class="spinner-border" role="status"></div>
            <div v-else>
                <div v-if="alert" :class="'alert alert-' + alert.type" role="alert">
                    {{alert.text}}
                </div>
                <div v-if="verifyPrompt" class="row">
                    <div class="col col-md-5">
                        <button type="button" class="btn btn-success" @click="acceptVerifyPromptClick()">Accept</button>
                    </div>
                    <div class="col col-md-5">
                        <button type="button" class="btn btn-danger" @click="denyVerifyPromptClick">Deny</button>
                    </div>
                </div>
                <form v-else>
                    <div class="row">
                        <input id="url-input" class="form-control" v-model="url">   
                    </div>
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
            verifyPrompt: true,
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
            console.log("URL: " + this.url)
        },
        denyVerifyPromptClick() {
            this.verifyPrompt = false
        },
        acceptVerifyPromptClick() {
            this.isLoading = true
            axios.get('/verify')
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
