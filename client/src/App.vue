<template>
    <div>
        <div id="navbar" class="ui fixed borderless huge inverted menu">
            <div class="header item"><b>VCD</b></div>
        </div>
        <div class="ui container">
            <LoadingSegment :isLoading="isQueryLoading">
                <div class="ui fluid action input">
                    <input type="text" v-model="url">
                    <button class="ui button">Submit</button>
                </div>
            </LoadingSegment>
            <LoadingSegment :isLoading="areCredsLoading">
                <div class="ui cards">
                    <CredCard v-for="cred in creds" :key="cred.name" :content="cred" />
                </div>
            </LoadingSegment>
        </div>
    </div>
</template>

<script>
import LoadingSegment from './components/LoadingSegment.vue'
import CredCard from './components/CredCard.vue'

const fs = require('fs')

export default {
    data() {
        return {
            isQueryLoading: false,
            areCredsLoading: false,
            creds: []
        }
    },
    components: {
        LoadingSegment, CredCard
    },
    created() {
        this.loadCreds()
    },
    methods: {
        loadCreds() {
            this.areCredsLoading = true
            fs.readFile('../wallet/verifiable-credentials.json', 'utf8', (err, data) => {
                this.areCredsLoading = false
                if (err) {
                    console.log(err)
                    return
                }

                console.log(data)
            })

            // this.creds = [
            //     {
            //         name: 'Student ID',
            //         issuer: 'University Issuer',
            //         issuer_did: 'university-issuer',
            //         fields: {
            //             'FirstName': 'Alice',
            //             'LastName': 'Student'
            //         }
            //     }
            // ]
        }
    }
}
</script>

<style>
body {
    padding-top: 4rem;
}

#navbar {
    background-color: #4a008a;
}

#alert-box {
    width: 100%;
}

#page-content {
    padding-top: 4rem;
    width: 100%;
}

.formInput {
    margin-bottom: 1rem;
}
</style>
