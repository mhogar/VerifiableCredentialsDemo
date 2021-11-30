<template>
<div>
    <div id="navbar" class="ui fixed borderless huge inverted menu">
        <div class="header item"><b>VCD</b></div>
    </div>
    <div class="ui container">
        <div class="ui center aligned basic segment">
            <div v-if="alert" id="alert-box" :class="'ui message ' + this.alert.type">
                <p class="center aligned">
                    {{alert.text}} 
                    <i class="close icon" @click="clearAlert"></i>
                </p>
            </div>
        </div>
        <div v-if="!prompt">
            <LoadingSegment :isLoading="isQueryLoading">
                <form>
                   <h2 class="ui center aligned header">Start a Query</h2>
                    <div class="ui fluid action input">
                        <input type="text" v-model="url">
                        <button type="submit" :class="'ui large icon button' + querySubmitDisabled" @click.prevent="querySubmit">
                            <i class="sign-in icon"></i>
                        </button>
                    </div>
                </form>
            </LoadingSegment>
            <LoadingSegment :isLoading="areCredsLoading">
                <div class="ui stackable three column grid">
                    <div class="column">
                        <CredCard v-for="(cred, issuer) in creds" :key="issuer" :issuer="issuer" :cred="cred" />
                    </div>
                </div>
            </LoadingSegment>
        </div>
        <Prompt v-else :prompt="prompt" :acceptCallback="acceptPromptCallback" :denyCallback="promptCallback" />
    </div>
</div>
</template>

<script>
import LoadingSegment from './components/LoadingSegment.vue'
import CredCard from './components/CredCard.vue'
import Prompt from './components/Prompt.vue'

import alertFactory from './common/alertFactory'
import http from './common/http'

export default {
    data() {
        return {
            isQueryLoading: false,
            areCredsLoading: false,
            alert: null,
            creds: {},
            url: 'http://localhost:8084/issue',
            prompt: null
        }
    },
    components: {
        LoadingSegment, CredCard, Prompt
    },
    created() {
        this.loadCreds()
    },
    computed: {
        querySubmitDisabled() {
            return this.url ? '' : ' disabled'
        }
    },
    methods: {
        setAlert(alert) {
            this.alert = alert
        },
        clearAlert() {
            this.alert = null
        },
        loadCreds() {
            this.areCredsLoading = true
            http.get('/creds')
            .then((res) => {
                if (res.data.error) {
                    this.setAlert(alertFactory.createErrorAlert(res.data.error))
                    return
                }

                this.creds = res.data
            })
            .catch((err) => {
                console.log(err)
                this.setAlert(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.areCredsLoading = false
            })
        },
        querySubmit() {
            this.isQueryLoading = true
            http.get('/query', {
                url: this.url
            })
            .then((res) => {
                if (res.data.error) {
                    this.setAlert(alertFactory.createErrorAlert('Query Failed: ' + res.data.error))
                    return
                }

                this.prompt = res.data
            })
            .catch((err) => {
                console.log(err)
                this.setAlert(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.isQueryLoading = false
            })
        },
        promptCallback(alert) {
            this.url = ''
            this.prompt = null

            if (alert) {
                this.setAlert(alert)
            }
        },
        acceptPromptCallback(alert, reloadCreds) {
            this.promptCallback(alert)

            if (reloadCreds) {
                this.loadCreds()
            }
        }
    }
}
</script>

<style>
body {
    padding-top: 4rem;
}

#navbar {
    background-color: #4a008a;
}

#alert-box {
    width: 100%;
}

#page-content {
    padding-top: 4rem;
    width: 100%;
}

.formInput {
    margin-bottom: 1rem;
}
</style>
