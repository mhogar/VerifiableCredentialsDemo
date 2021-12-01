<template>
<div>
    <div id="type-header" class="ui basic segment">
        <h1 class="ui center aligned header">
            {{typeTitle}}
            <div class="sub header">{{typeDescription}}</div>
        </h1>
    </div>
    <Alert ref="alert" />
    <LoadingSegment :isLoading="isPromptLoading">
        <div v-if="!showForm">
            <h3 class="ui header">{{entityType}} Information:</h3>
            <div  class="ui fluid raised card">
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
                    <button type="button" :class="'ui primary button' + acceptButtonDisabled" @click="acceptButtonClicked">Accept</button>
                    <button type="button" class="ui button" @click="denyButtonClicked">Deny</button>
                </div>
            </div>
        </div>
        <Form v-else :url="prompt.service_url" :fields="prompt.fields" :submitCallback="submitFormCallback" />
    </LoadingSegment>
    <LoadingSegment v-if="hasIssuer && cred" :isLoading="isCredLoading">
        <h3 class="ui header">Applicable Credentials:</h3>
        <div class="ui stackable three column grid">
            <div class="column">
                <CredCard :cred="cred" />
            </div>
        </div>
    </LoadingSegment>
</div>
</template>

<script>
import Alert from './Alert.vue'
import LoadingSegment from './LoadingSegment.vue'
import CredCard from './CredCard.vue'
import Form from './Form.vue'

import alertFactory from '../common/alertFactory'
import http from '../common/http'

export default {
    data() {
        return {
            isPromptLoading: false,
            isCredLoading: false,
            showForm: false,
            cred: null
        }
    },
    components: {
        Alert, LoadingSegment, CredCard, Form
    },
    props: {
        prompt: Object,
        acceptCallback: Function,
        denyCallback: Function
    },
    created() {
        if (this.hasIssuer) {
            this.loadCred()
        }
    },
    computed: {
        hasIssuer() {
            const type = this.prompt.type
            return type === 'verify' || type === 'iss:cred'
        },
        acceptButtonDisabled() {
            return this.hasIssuer && !this.cred ? ' disabled' : ''
        },
        typeTitle() {
            switch (this.prompt.type) {
                case 'verify':
                    return 'Verify Request'
                case 'iss:form':
                case 'iss:cred':
                    return 'Issue Request'
                default:
                    return ''
            }
        },
        typeDescription() {
            switch (this.prompt.type) {
                case 'verify':
                    return 'Verify a credential.'
                case 'iss:form':
                    return 'Create a credential by filling out a form.'
                case 'iss:cred':
                    return 'Create a new credential from an existing one.'
                default:
                    return ''
            }
        },
        entityType() {
            switch (this.prompt.type) {
                case 'verify':
                    return 'Verifier'
                case 'iss:form':
                case 'iss:cred':
                    return 'Issuer'
                default:
                    return ''
            }
        },
        trustedByVerifierIcon() {
            return this.prompt.trusted_by_issuer ? 'check circle green icon' : 'close red icon'
        }
    },
    methods: {
        setAlert(alert) {
            this.$refs.alert.setAlert(alert)
        },
        acceptButtonClicked() {
            this.setAlert(null)

           switch (this.prompt.type) {
                case 'verify':
                    this.verify()
                    return
                case 'iss:form':
                    this.showForm = true
                    return
                case 'iss:cred':
                    this.issueCred()
                    return
            }
        },
        denyButtonClicked() {
            this.denyCallback(alertFactory.createWarningAlert('Request denied.'))
        },
        submitFormCallback(alert, reloadCreds) {
            this.acceptCallback(alert, reloadCreds)
        },
        loadCred() {
            this.isCredLoading = true
            http.get('/cred', {
                id: this.prompt.issuer
            })
            .then((res) => {
                if (!res.data.issuer) {
                    this.setAlert(alertFactory.createWarningAlert('No credentials for target issuer found.'))
                    return
                }
                this.cred = res.data
            })
            .catch((err) => {
                console.log(err)
                this.setAlert(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.isCredLoading = false
            })
        },
        verify() {
            this.isPromptLoading = true
            http.post('/verify', {
                service_url: this.prompt.service_url,
                credential_id: this.prompt.issuer
            })
            .then((res) => {
                if (res.data.error) {
                    this.setAlert(alertFactory.createErrorAlert('Request Failed: ' + res.data.error))
                    return
                }

                this.acceptCallback(alertFactory.createSuccessAlert('Verified!'))
            })
            .catch((err) => {
                console.log(err)
                this.setAlert(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.isPromptLoading = false
            })
        },
        issueCred() {
            this.isPromptLoading = true
            http.post('/issue', {
                service_url: this.prompt.service_url,
                type: 'iss:cred',
                credential_id: this.prompt.issuer
            })
            .then((res) => {
                if (res.data.error) {
                    this.setAlert(alertFactory.createErrorAlert('Create VC Failed: ' + res.data.error))
                    return
                }

                this.acceptCallback(alertFactory.createSuccessAlert('Created Credential!'), true)
            })
            .catch((err) => {
                console.log(err)
                this.setAlert(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.isPromptLoading = false
            })
        }
    }
}
</script>

<style scoped>
#type-header {
    padding-bottom: 2rem;
}
</style>
