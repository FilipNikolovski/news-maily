/**
 * Created by filip on 25.7.15.
 */
var Template = function () {

    var url = url_base + '/api/templates';

    this.all = function (paginate, perPage, page) {
        return $.ajax({
            type: 'GET',
            url: url,
            data: {
                paginate: paginate,
                per_page: perPage,
                page: page
            },
            dataType: 'json'
        });
    };

    this.get = function (id) {
        return $.ajax({
            type: 'GET',
            url: url + '/' + id,
            dataType: 'json'
        });
    };

    this.delete = function (id) {
        return $.ajax({
            type: 'DELETE',
            url: url + '/' + id,
            dataType: 'json'
        });
    };

    this.create = function (data) {
        return $.ajax({
            type: 'POST',
            url: url,
            data: {
                name: data.name,
                content: data.content
            },
            dataType: 'json'
        });
    }
};

module.exports = Template;