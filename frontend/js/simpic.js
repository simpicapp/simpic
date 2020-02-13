import Vue from './vue'
import timeline from './timeline'

const app = new Vue({
    components: {
        timeline
    },
    el: '#main',
    template: `
    <main>
        <timeline/>
    </main>
    `
});
