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
                        <i class="close icon" @click="clearAlert()"></i>
                    </p>
                </div>
            </div>
            <LoadingSegment :isLoading="isQueryLoading">
                <div class="ui fluid action input">
                    <input type="text" v-model="url">
                    <button class="ui button">Submit</button>
                </div>
            </LoadingSegment>
            <LoadingSegment :isLoading="areCredsLoading">
                <div class="ui cards">
                    <CredCard v-for="(cred, issuer) in creds" :key="issuer" :issuer="issuer" :cred="cred" />
                </div>
            </LoadingSegment>
        </div>
    </div>
</template>

<script>
import LoadingSegment from './components/LoadingSegment.vue'
import CredCard from './components/CredCard.vue'

import http from './common/http'

export default {
    data() {
        return {
            alert: null,
            isQueryLoading: false,
            areCredsLoading: false,
            creds: {}
        }
    },
    components: {
        LoadingSegment, CredCard
    },
    created() {
        this.loadCreds()
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
        loadCreds() {
            this.areCredsLoading = true
            http.get('/creds', {
                validateStatus: (status) => status < 500
            })
            .then((res) => {
                if (res.data.error) {
                    this.setAlert("negative", res.data.error)
                    return
                }

                this.creds = res.data
            })
            .catch((err) => {
                console.log(err)
                this.setAlert("negative", "An internal error occurred. Try again later.")
            })
            .then(() => {
                this.areCredsLoading = false
            })
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
