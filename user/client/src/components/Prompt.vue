<template>
<div>
    <div id="type-header" class="ui basic segment">
        <h1 class="ui center aligned header">
            {{typeTitle}}
            <div class="sub header">{{typeDescription}}</div>
        </h1>
    </div>
    <div class="ui stackable centered three column grid">
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
                    <button type="button" class="ui positive button" @click="acceptButtonClicked()">Accept</button>
                    <div class="or"></div>
                    <button type="button" class="ui negative button" @click="denyButtonClicked()">Deny</button>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import alertFactory from '../common/alertFactory'

export default {
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
            const type = this.prompt.type
            switch (type) {
                case 'verify':
                    return 'Verify Request'
                case 'iss:form':
                case 'iss:cred':
                    return 'Issue Request'
            }
            
            return 'Unknown Type'
        },
        typeDescription() {
            const type = this.prompt.type
            switch (type) {
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
            this.acceptCallback(alertFactory.createSuccessAlert('Request accepted!'))
        },
        denyButtonClicked() {
            this.denyCallback(alertFactory.createWarningAlert('Request denied.'))
        },
    }
}
</script>

<style scoped>
#type-header {
    padding-bottom: 2rem;
}
</style>
