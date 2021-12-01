<template>
<div>
    <div v-if="alert" class="ui center aligned basic segment">
        <div :class="'ui message ' + alert.type">
            <p class="center aligned">
                {{alert.text}} 
                <i class="close icon" @click="clearAlert"></i>
            </p>
        </div>
    </div>
</div>
</template>

<script>
export default {
    data() {
        return {
            alert: null
        }
    },
    methods: {
        setAlert(alert) {
            this.alert = alert
        },
        clearAlert() {
            this.alert = null
        }
    }
}
</script>

<style scoped>
.ui.message {
    width: 100%;
}
</style>
