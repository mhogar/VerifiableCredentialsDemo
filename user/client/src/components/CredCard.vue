<template>
<div class="ui fluid raised card">
    <div class="content">
        <div class="header">
            {{cred.cred_type}}
        </div>
        <div class="meta">
            Issued by: {{cred.issuer.did}}
        </div>
    </div>
    <div class="content">
        <div class="description">
            <p v-for="(value, key) in cred.credentials" :key="key">
                <b>{{key}}: </b> {{value}}
            </p>
        </div>
    </div>
</div>
</template>

<script>
export default {
    props: {
        cred: Object
    }
}
</script>

<style scoped>
.ui.card {
    background-color: aliceblue;
}
</style>
