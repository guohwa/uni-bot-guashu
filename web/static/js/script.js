$(document).ready(function() {
    toastr.options = {
        "closeButton": true,
        "progressBar": true,
        "positionClass": "toast-bottom-right"
    };
    $('.select2').select2({minimumResultsForSearch:-1});
    $('.bi-eye-fill').click(function() {
        var field = $(this).prev('.form-control');
        if (field.attr('type') == 'text') {
            field.attr('type', 'password');
        } else {
            field.attr('type', 'text');
        }
    });
    $('.ajax-link').click(function(evt) {
        var that = $(this),
            action = that.attr('href'),
            timeout = 3000,
            loading = $('.loading');
        evt.preventDefault();
        loading.show();
        $.ajax({
            url: action,
            success: function(data) {
                if (data.code) {
                    toastr.error(data.msg);
                } else {
                    toastr.success(data.msg);
                    setTimeout(function() {
                        location.href = data.url;
                    }, timeout);
                }
            },
            complete: function() {
                loading.hide();
            }
        });
    });
    $('.ajax-form').submit(function(evt) {
        var that = $(this),
            action = that.attr('action'),
            method = that.attr('method'),
            timeout = 3000,
            loading = $('.loading');
        evt.preventDefault();
        loading.show();
        $.ajax({
            url: action,
            method: method,
            data: that.serialize(),
            success: function(data) {
                if (data.code) {
                    toastr.error(data.msg);
                    that.find('.captcha').click();
                } else {
                    toastr.success(data.msg);
                    setTimeout(function() {
                        location.href = data.url;
                    }, timeout);
                }
            },
            complete: function() {
                loading.hide();
            }
        });
    });

    $('.btn-close').click(function(evt) {
        var that = $(this),
            href = that.attr('href'),
            loading = $('.loading'),
            timeout = 3000;
        var data = {
            symbol: that.data('symbol'),
            customer: that.data('customer'),
            positionSide: that.data('positionside'),
            positionAmt: that.data('positionamt')
        };
        evt.preventDefault();
        loading.show();
        $.ajax({
            url: href,
            method: 'POST',
            data: data,
            success: function(data) {
                if (data.code) {
                    toastr.error(data.msg);
                } else {
                    toastr.success(data.msg);
                    setTimeout(function() {
                        location.href = data.url;
                    }, timeout);
                }
            },
            complete: function() {
                loading.hide();
            }
        });
    });
    $('.btn-cancel').click(function(evt) {
        var that = $(this),
            href = that.attr('href'),
            loading = $('.loading'),
            timeout = 3000;
        var data = {
            symbol: that.data('symbol'),
            customer: that.data('customer'),
            orderId: that.data('orderid')
        };
        evt.preventDefault();
        loading.show();
        $.ajax({
            url: href,
            method: 'POST',
            data: data,
            success: function(data) {
                if (data.code) {
                    toastr.error(data.msg);
                } else {
                    toastr.success(data.msg);
                    setTimeout(function() {
                        location.href = data.url;
                    }, timeout);
                }
            },
            complete: function() {
                loading.hide();
            }
        });
    });

    $('.btn-order').click(function(evt) {
        var that = $(this),
            loading = $('.loading'),
            form = $(that.parents("form")),
            url = form.attr('action'),
            method = form.attr('method'),
            timeout = 3000,
            side = that.data('value'),
            data = form.serialize() + "&" + $.param({side});
        evt.preventDefault();
        loading.show();
        $.ajax({
            url,
            method,
            data,
            success: function(data) {
                if (data.code) {
                    toastr.error(data.msg);
                    that.find('.captcha').click();
                } else {
                    toastr.success(data.msg);
                    setTimeout(function() {
                        location.href = data.url;
                    }, timeout);
                }
            },
            complete: function() {
                loading.hide();
            }
        });
    });
});
