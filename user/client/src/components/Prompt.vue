<template>
<LoadingSegment :isLoading="isLoading">
    <div id="type-header" class="ui basic segment">
        <h1 class="ui center aligned header">
            {{typeTitle}}
            <div class="sub header">{{typeDescription}}</div>
        </h1>
    </div>
    <div v-if="!showForm" class="ui stackable centered three column grid">
        <div class="ui raised card">
            <div class="content">
                <div class="header">{{prompt.name}}</div>
                <div class="meta">
                    <span>{{prompt.domain}}</span>
                </div>
                <h4 v-if="showTrusted" class="ui sub header">
                    Trusted By Issuer:
                    <i :class="trustedByVerifierIcon"></i>
                </h4>
            </div>
            <div class="content">
                <div class="description">
                    <p>{{prompt.purpose}}</p>
                </div>
            </div>
            <div class="extra content">
                <div class="ui buttons">
                    <button type="button" class="ui positive button" @click="acceptButtonClicked">Accept</button>
                    <div class="or"></div>
                    <button type="button" class="ui negative button" @click="denyButtonClicked">Deny</button>
                </div>
            </div>
        </div>
    </div>
    <Form v-else :url="prompt.service_url" :fields="prompt.fields" :submitCallback="submitFormCallback" />
</LoadingSegment>
</template>

<script>
import LoadingSegment from './LoadingSegment.vue'
import Form from './Form.vue'

import alertFactory from '../common/alertFactory'
import http from '../common/http'

export default {
    data() {
        return {
            isLoading: false,
            showForm: false
        }
    },
    components: {
        LoadingSegment, Form
    },
    props: {
        prompt: Object,
        acceptCallback: Function,
        denyCallback: Function
    },
    computed: {
        showTrusted() {
            const type = this.prompt.type
            return type === 'verify' || type === 'iss:cred'
        },
        typeTitle() {
            switch (this.prompt.type) {
                case 'verify':
                    return 'Verify Request'
                case 'iss:form':
                case 'iss:cred':
                    return 'Issue Request'
            }
            
            return 'Unknown Type'
        },
        typeDescription() {
            switch (this.prompt.type) {
                case 'verify':
                    return 'Verify a credential.'
                case 'iss:form':
                    return 'Create a credential by filling out a form.'
                case 'iss:cred':
                    return 'Create a new credential from an existing one.'
            }
            
            return ''
        },
        trustedByVerifierIcon() {
            return this.prompt.trusted_by_issuer ? 'check circle green icon' : 'close red icon'
        }
    },
    methods: {
        acceptButtonClicked() {
           switch (this.prompt.type) {
                case 'verify':
                    this.verify()
                    return
                case 'iss:form':
                    this.showForm = true
                    return
                case 'iss:cred':
                    this.verify()
                    return
            }

            this.acceptCallback()
        },
        denyButtonClicked() {
            this.denyCallback(alertFactory.createWarningAlert('Request denied.'))
        },
        submitFormCallback(alert, reloadCreds) {
            this.acceptCallback(alert, reloadCreds)
        },
        verify() {
            this.isLoading = true
            http.post('/verify', {
                service_url: this.prompt.service_url,
                credential_id: this.prompt.issuer
            })
            .then((res) => {
                if (res.data.error) {
                    this.acceptCallback(alertFactory.createErrorAlert('Request Failed: ' + res.data.error))
                    return
                }

                this.acceptCallback(alertFactory.createSuccessAlert('Verified!'))
            })
            .catch((err) => {
                console.log(err)
                this.acceptCallback(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.isLoading = false
            })
        },
        issueCred() {
        }
    }
}
</script>

<style scoped>
#type-header {
    padding-bottom: 2rem;
}
</style>
