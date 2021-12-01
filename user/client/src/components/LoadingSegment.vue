<template>
<div class="ui basic segment">
    <div v-if="isLoading" class="ui active loader"></div>
    <slot v-else />
</div>
</template>

<script>
export default {
    props: {
        isLoading: Boolean
    }
}
</script>
