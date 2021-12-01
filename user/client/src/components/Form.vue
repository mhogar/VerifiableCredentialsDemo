<template>
<div>
    <Alert ref="alert" />
    <LoadingSegment :isLoading="isLoading">
        <form class="ui form">
            <div class="field" v-for="field in fields" :key="field.name">
                <label>{{field.name}}</label>
                <input :type="field.type" v-model="values[field.name]">
            </div>
            <button type="submit" :class="'ui primary button' + submitDisabledClass" @click.prevent="submit">Submit</button>
            <button class="ui button" @click.prevent="cancel">Cancel</button>
        </form>
    </LoadingSegment>
</div>
</template>

<script>
import Alert from './Alert.vue'
import LoadingSegment from './LoadingSegment.vue'

import alertFactory from '../common/alertFactory'
import http from '../common/http'

export default {
    data() {
        return {
            isLoading: false,
            values: {}
        }
    },
    components: {
        Alert, LoadingSegment
    },
    props: {
        url: String,
        fields: Array,
        submitCallback: Function
    },
    created() {
        this.fields.forEach(field => this.values[field.name] = '');
    },
    computed: {
        canSubmit() {
            return Object.values(this.values).find(val => !val) == null
        },
        submitDisabledClass() {
            return this.canSubmit ? '' : ' disabled'
        }
    },
    methods: {
        setAlert(alert) {
            this.$refs.alert.setAlert(alert)
        },
        cancel() {
            this.submitCallback(alertFactory.createWarningAlert('Issue request canceled.'))
        },
        submit() {
            if (!this.canSubmit) {
                return
            }

            this.setAlert(null)

            this.isLoading = true
            http.post('/issue', {
                service_url: this.url,
                type: 'iss:form',
                fields: this.values
            })
            .then((res) => {
                if (res.data.error) {
                    this.setAlert(alertFactory.createErrorAlert('Create VC Failed: ' + res.data.error))
                    return
                }

                this.submitCallback(alertFactory.createSuccessAlert('Created Credential!'), true)
            })
            .catch((err) => {
                console.log(err)
                this.setAlert(alertFactory.createInternalErrorAlert())
            })
            .then(() => {
                this.isLoading = false
            })
        }
    }
}
</script>

<style scoped>
.ui.labeled.input {
    margin-bottom: 1rem;
}

.ui.form {
    margin-bottom: 1rem;
}
</style>
