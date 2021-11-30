<template>
    <div class="ui raised purple card">
        <div class="content">
            <div class="header">
                {{content.name}}
            </div>
            <div class="meta">
                {{content.issuer}} ({{content.issuer_did}})
            </div>
        </div>
        <div class="content">
            <div class="description">
                <p v-for="(value, key) in content.fields" :key="key">
                    <b>{{key}}: </b> {{value}}
                </p>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        content: Object
    }
}
</script>
