/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');
require('sweetalert');

var React = require('react');
var List = require('../../entities/list.js');
var l = new List();

var DeleteButton = React.createClass({
    handleSubmit: function (e) {
        e.preventDefault();
        swal({
                title: "Are you sure?",
                text: "You will not be able to recover this subscriber!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "Yes, delete it!",
                closeOnConfirm: false
            },
            function () {
                l.deleteSubscriber(this.props.sid)
                    .done(function () {
                        swal({
                            title: "Success",
                            text: "The subscriber was successfully removed!",
                            type: "success"
                        }, function () {
                            location.reload();
                        });
                    })
                    .fail(function () {
                        swal('Could not delete', 'Could not delete the subscriber. Try again.', 'error');
                    });
            }.bind(this));
    },
    render: function () {
        return (
            <form onSubmit={this.handleSubmit}>
                <input type="hidden" name="_method" value="DELETE"/>
                <button type="submit"><span className="glyphicon glyphicon-trash"></span></button>
            </form>
        );
    }
});

var SubscriberRow = React.createClass({
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td>{this.props.data.email}</td>
                <td>
                    <DeleteButton sid={this.props.data.id}/>
                </td>
            </tr>
        );
    }
});

var SubscribersTable = React.createClass({
    getInitialState: function () {
        return {subscribers: {data: []}};
    },
    componentDidMount: function () {
        l.getSubscribers(this.props.listId, true, 10, 1).done(function (response) {
            this.setState({subscribers: response});
            $('.pagination').bootpag({
                total: response.last_page,
                page: response.current_page,
                maxVisible: 5
            }).on("page", function (event, num) {
                l.getSubscribers(this.props.listId, true, 10, num).done(function (response) {
                    this.setState({subscribers: response});
                    $('.pagination').bootpag({page: response.current_page});
                }.bind(this));
            }.bind(this));
        }.bind(this));
    },
    render: function () {
        var rows = function (data) {
            return <SubscriberRow key={data.id} data={data}/>
        }.bind(this);
        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Subscriber name</th>
                        <th>Email</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.subscribers.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
});

module.exports = SubscribersTable;
