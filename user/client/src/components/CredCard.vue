<template>
<div class="ui fluid purple card">
    <div class="content">
        <!-- <div class="header">
            {{content.name}}
        </div>
        <div class="meta">
            {{content.issuer}} ({{content.issuer_did}})
        </div> -->
        <div class="header">
            {{issuer}}
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
        issuer: String,
        cred: Object
    }
}
</script>
