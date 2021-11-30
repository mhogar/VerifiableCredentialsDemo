<template>
    <LoadingSegment :isLoading="isLoading">
        <form class="ui form">
            <div class="field" v-for="field in fields" :key="field.name">
                <label>{{field.name}}</label>
                <input :type="field.type" v-model="values[field.name]">
            </div>
        </form>
        <button :class="'ui primary button' + submitDisabled" @click.prevent="submit">Submit</button>
        <button class="ui button" @click.prevent="cancel">Cancel</button>
    </LoadingSegment>
</template>

<script>
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
        LoadingSegment
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
        submitDisabled() {
            return Object.values(this.values).find(val => !val) != null ? ' disabled' : ''
        }
    },
    methods: {
        cancel() {
            this.submitCallback(alertFactory.createWarningAlert('Issue request canceled.'))
        },
        submit() {
            this.isLoading = true
            http.post('issue', {
                service_url: this.url,
                fields: this.values
            })
            .then((res) => {
                if (res.data.error) {
                    this.submitCallback(alertFactory.createErrorAlert('Create VC Failed: ' + res.data.error))
                    return
                }

                this.submitCallback(alertFactory.createSuccessAlert('Verified!'))
            })
            .catch((err) => {
                console.log(err)
                this.submitCallback(alertFactory.createInternalErrorAlert())
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
