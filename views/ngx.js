function ajax(option) {
    return new Promise((resolve, reject) => {
      let xhr = new XMLHttpRequest()
      xhr.open(option.method || 'GET', option.url)
      xhr.onload = () => {
        // let data = JSON.parse(xhr.responseText)
        resolve(xhr.responseText)
      }
      xhr.onerror = error => {
        reject(error)
      }
      xhr.send(option.data ? JSON.stringify(option.data) : null)
    })
}

new Vue({
    el: '#ngxman',
    data: {
        sites: [
            'a.cloudhua.com',
            'b.cloudhua.com'
        ],
        sites2: [
            {
                value: 'a.cloudhua.com',
                label: 'a.cloudhua.com',
            },
            {
                value: 'b.cloudhua.com',
                label: 'b.cloudhua.com',
            }
        ],
        currentSite: '',
        myCodeMirror: null,
        theme: 'light'
    },
    mounted() {
        this.myCodeMirror = CodeMirror.fromTextArea(document.getElementById('ngxcfg'), {
            lineNumbers: true,
            theme: 'night',
            lineWrapping: true,
            mode: 'nginx',
        });
        this.myCodeMirror.setSize("100%", "100%")
    },
    methods: {
        showsite: function(site) {
            this.currentSite = site
            this.myCodeMirror.setValue(site)
        },
    },
    components: {
        'site-item': {
            template: `
            <menu-item :name=site>
            <icon type="link"></icon>
            {{ site }}
            </menu-item>
            `,
            props: ['site'],
        },
    }
})

function tranSite(site) {
    return site.replace(/\./g, '-')
}

function isValidDomain(v) {
    if (typeof v !== 'string') return false

    var parts = v.split('.')
    if (parts.length <= 2) return false

    var tld = parts.pop()
    var tldRegex = /^[a-zA-Z0-9]+$/gi

    if (!tldRegex.test(tld)) return false

    var isValid = parts.every(function(host) {
        var hostRegex = /^(?!:\/\/)([a-zA-Z0-9]+|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])$/gi;
        return hostRegex.test(host)
    })

    return isValid
}