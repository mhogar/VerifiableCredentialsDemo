export default {
    createAlert(type, text) {
        return {
            type: type,
            text: text
        }
    },
    createSuccessAlert(text) {
        return this.createAlert('positive', text)
    },
    createWarningAlert(text) {
        return this.createAlert('warning', text)
    },
    createErrorAlert(text) {
        return this.createAlert('negative', text)
    },
    createInternalErrorAlert() {
        return this.createErrorAlert('An internal error occurred. Try again later.')
    }
}
