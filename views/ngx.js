
var sites = [
];

function alert(data) {
    $('#alert').html(data)
    $('.alert-ui.ui.tiny.modal').modal('show');
}

function createSite(e) {
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
        $('#ngx-cfg-view').html(data)
        sites.push({title: site})
        $('#sites').append('<div class="item"><i class="site linkify icon" style="float: left;">' + site + '</i></div>')
    })
}

function saveCfg(newname) {
    var site = $('#currentsite').html()
    var siteData = $('#ngx-cfg-view').html()
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
                $('#ngx-cfg-view').html(siteData)
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

function showCfg(site) {
    $.get('sitecontent?site=' + encodeURIComponent(site), function(data) {
        $('#ngx-cfg-view').html(data['content'])
    }).fail(function() {
        alert('fail')
    })
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

$(document).ready(function() {

    $.each($('#sites .site').toArray(), function(i, v) {
        sites[i] = {title: v['innerText']}
    })

    $('.ui.search').search({
        source: sites,
        onSelect: function(result, resp) {
            var site = result["title"]
            showCfg(site)
        }
    });

    $('#sites .item').on('click', function() {
            $(this).addClass('active').siblings().removeClass('active');
            var site = $(this).children(".site")[0].innerHTML
            $('#currentsite').html(site)
            showCfg(site)
        }
    );

    dbclickRename('#sites .item .site')

    $('#create-site').on('click', createSite)
    $('#save-cfg').on('click', saveCfg)
    $('#test-cfg').on('click', testCfg)
    $('#reload-cfg').on('click', reloadCfg)
})
