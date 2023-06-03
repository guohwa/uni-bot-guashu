$(document).ready(function() {
    $(".pagger").each(function(idx, elem) {
        let pagger = $(elem),
            limit = pagger.data('limit'),
            page = pagger.data('page'),
            href = pagger.data('href'),
            count = pagger.data('count');
        let p = page - 1,
            c = Math.ceil(count / limit),
            b = -1,
            e = -1;
        if (p - 2 >= 0) {
            b = p - 2; if (p + 2 > c -1) {e = c - 1;} else {e = p + 2}
        } else {
            b = 0; e = 4;
        } 
        e = e + 1; e = e > c ? c : e;

        let pagination = $('<ul class="pagination pagination-sm"></ul>');
        let previous = $(`<li class="page-item previous"><a class="page-link" aria-label="Previous"><span aria-hidden="true">&laquo;</span></a></li>`);
        let next = $(`<li class="page-item next"><a class="page-link" aria-label="Next"><span aria-hidden="true">&raquo;</span></a></li>`);
        pagger.append(pagination);
        pagination.append(previous);
        pagination.append(next);

        if (p == 0 || count <= 0) {
            previous.addClass('disabled');
        }
        if (p > 0 && count > 0) {
            previous.children('.page-link').attr('href', `${href}?page=${p}`);
        }
        if (b == 0 && e == 0) {
            next.before(`<li class="page-item active"><a class="page-link" href="${href}?page=1">1</a></li>`);
        }
        if (b > 0) {
            next.before(`<li class="page-item"><a class="page-link" href="${href}?page=1">1</a></li>`);
        }
        if (b > 1) {
            next.before(`<li class="page-item disabled"><a class="page-link" aria-label="">...</a></li>`);
        }
        for (let i = b; i < e; i++) {
            let active = p == i ? ' active' : '';
            let aria = p == i ? ' aria-current="page"' : '';
            next.before(`<li class="page-item${active}"${aria}>
                <a class="page-link" href="${href}?page=${i + 1}">${i + 1}</a>
            </li>`);
        }
        if (e < c - 1) {
            next.before(`<li class="page-item disabled"><a class="page-link" aria-label="">...</a></li>`);
        }
        if (e < c) {
            next.before(`<li class="page-item"><a class="page-link" href="${href}?page=${c}">${c}</a></li>`);
        }
        if (p === c - 1 || count <= 0) {
            next.addClass('disabled');
        }
        if (p < c - 1 && count > 0) {
            next.children('.page-link').attr('href', `${href}?page=${p+2}`);
        }
    });
});
