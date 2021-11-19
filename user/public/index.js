const app = {
    template: `
        <div class="container">
            <div v-if="verifyPrompt" class="row justify-content-center">
                <div class="col col-md-1">
                    <button type="button" class="btn btn-success" @click="verifyPromptClick(true)">Accept</button>
                </div>
                <div class="col col-md-1">
                    <button type="button" class="btn btn-danger" @click="verifyPromptClick(false)">Deny</button>
                </div>
            </div>
            <form v-else class="row justify-content-center">
                <div class="col col-md-5">
                    <input id="url-input" class="form-control" v-model="url">
                    <button type="submit" class="btn btn-primary" @click.prevent="submitPresReq">Submit</button>
                </div>
            </form>
        </div>
    `,
    data() {
        return {
            verifyPrompt: true,
            url: ""
        }
    },
    methods: {
        submitPresReq() {
            console.log("URL: " + this.url)
        },
        verifyPromptClick(result) {
            console.log("Verify: " + result)
            this.verifyPrompt = false
        }
    }
};

Vue.createApp(app).mount('#app')
