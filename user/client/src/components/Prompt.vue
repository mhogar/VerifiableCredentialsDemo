<template>
<LoadingSegment :isLoading="isLoading">
    <div id="type-header" class="ui basic segment">
        <h1 class="ui center aligned header">
            {{typeTitle}}
            <div class="sub header">{{typeDescription}}</div>
        </h1>
    </div>
    <div v-if="!showForm" class="ui fluid raised card">
        <div class="content">
            <div class="header">{{prompt.name}} - {{prompt.domain}}</div>
            <div class="meta">
                <p>{{prompt.did}}</p>
            </div>
            
        </div>
        <div class="content">
            <div class="description">
                <p v-if="hasIssuer"><b>Target Issuer: </b>{{prompt.issuer}}</p>
                <p><b>Credential Type: </b>{{prompt.cred_type}}</p>
                <p><b>Description: </b>{{prompt.description}}</p>
            </div>
            <h4 v-if="hasIssuer" class="ui sub header">
                Trusted By Target Issuer:
                <i :class="trustedByVerifierIcon"></i>
            </h4>
        </div>
        <div class="extra content">
            <button type="button" class="ui primary button" @click="acceptButtonClicked">Accept</button>
            <button type="button" class="ui button" @click="denyButtonClicked">Deny</button>
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
        hasIssuer() {
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
