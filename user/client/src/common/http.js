const axios = require('axios')

const http = axios.create({
    baseURL: "http://localhost:8082",
    headers: {
        "Content-type": "application/json"
    },
    validateStatus(status) {
        return status >= 200 && status < 500
    }
})

export default {
    get(url, params) {
        return http.get(url, {
            params: params
        })
    },
    post(url, data) {
        return http.post(url, data)
    }
}
