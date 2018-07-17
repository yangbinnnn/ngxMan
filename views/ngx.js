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
        this.myCodeMirror.setSize("100%", "700px")
    },
    methods: {
        showcfg: function(site) {
            this.currentSite = site
            this.myCodeMirror.setValue(site)
        },
    },
    components: {
        'site-item': {
            template: `
            <div> 
            <button v-on:click="$emit('showcfg', site)">{{ site }}</button>
            <i class="linkify icon" style="float: left;"></i>
            </div>
            `,
            props: ['site'],
        },
    }
})

function alert(data) {
    $('#alert').html(data)
    $('.alert-ui.ui.tiny.modal').modal('show');
}

function createSite(e) {
    alert('createSite')
    return
    var site = $('#new-site').val()
    if (!isValidDomain(site)) {
        alert("无效的域名\n示例:abc.xyz.com")
        return
    }
    exsit = false
    $.each(sites, function(idx, item) {
        if (item['title'] === site) {
            exsit = true
        }
    })
    if (exsit) {
        alert("站点已存在")
        return
    }
    $.get('createSite?site=' + encodeURIComponent(site), function(data) {
        $('#currentsite').html(site)
        myCodeMirror.setValue(data)
        sites.push({title: site})
        $('#sites').prepend('<div id="' + tranSite(site) + '" class="item"><i class="site linkify icon" style="float: left;">' + site + '</i></div>')
        $('#'+tranSite(site)).addClass('active').siblings().removeClass('active');
    })
}

function saveCfg(newname) {
    alert('saveCfg')
    var site = $('#currentsite').text()
    var siteData = myCodeMirror.getValue()
    if (site === '') {
        alert("未选择站点")
        return false
    }
    if (siteData.length === 0) {
        alert("站点内容为空")
        return false
    }
    if (!isValidDomain(site)) {
        alert("无效的域名\n示例:abc.xyz.com")
        return
    }

    var ok = true
    var verifyName = site
    if (typeof newname === 'string') {
        verifyName = newname
    }
    $('#confirm-site-verify').text(verifyName)
    $('.confirm-ui.ui.mini.modal').modal({
        onApprove: function() {
            if ($("#confirm-site").val() != verifyName) {
                alert("输入站点不匹配")
                ok = false
                return
            }

            url = 'savesite?site=' + encodeURIComponent(site)
            if (typeof newname  === 'string') {
                if (!isValidDomain(newname)) {
                    alert("无效的域名\n示例:abc.xyz.com")
                    return
                }
                url = url + '&rename=' + encodeURIComponent(newname)
                siteData = siteData.replace(new RegExp(site, 'g'), newname)
                myCodeMirror.setValue(siteData)
            }
            $.post(url, siteData, function() {
                alert("保存成功，测试并加载配置后生效")
            }).fail(function(data) {
                ok = false
                alert(data['responseText'])
            })
        }
    }).modal('show')
    return ok
}

function testCfg() {
    $.get('testsite', function(data) {
        alert(data)
    }).fail(function(data) {
        alert(data)
    })
}

function reloadCfg() {
    $.get('reloadsite', function(data) {
        alert(data)
    }).fail(function(data) {
        alert(data)
    })
}

function showCfg1(site) {
    // data = ajax({
    //     url: 'sitecontent?site=' + encodeURIComponent(site)
    // })
    data = 'hello ' + site
    if (data != null) {
        this.currentSite = site
        // this.myCodeMirror.setValue(site)
    }
    console.log(this.currentSite)
    console.log(this.sites)
    console.log(data)
    return

    $.get('sitecontent?site=' + encodeURIComponent(site), function(data) {
        $('#currentsite').html(site)
        $('#'+tranSite(site)).addClass('active').siblings().removeClass('active');
        myCodeMirror.setValue(data['content'])
    }).fail(function() {
        alert('获取内容失败')
    })
}

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

function renameSite(oldsite, newsite) {
        if (oldsite === newsite) {
            return false
        }
        if (!isValidDomain(newsite)) {
            alert("无效的域名\n示例:abc.xyz.com")
            return false
        }
        exsit = false
        $.each(sites, function(idx, item) {
            if (item['title'] === newsite) {
                exsit = true
            }
        })
        if (exsit) {
            alert("站点已存在")
            return false
        }
        var ok = saveCfg(newsite)
        if (ok) {
            $.each(sites, function(idx, item) {
                if (item['title'] === oldsite) {
                    item['title'] = newsite
                }
            })
            $('#currentsite').html(newsite)
            return true
        }
        return false
}

function dbclickRename(selector) {
    var orival
    var ok
    $(selector).on('dblclick', function() {
        orival = $(this).text()
        var item = this
        $('#rename').text(orival)
        $('.rename-ui.ui.mini.modal').modal({
            onApprove : function() {
                var newsite = $('#rename').text()
                ok = renameSite(orival, newsite)
                if (ok) {
                    $(item).text(newsite)
                }
            }
        }).modal('show');
    })
}
