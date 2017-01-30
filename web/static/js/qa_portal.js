var qa_portal = {
	initChallenges: function() {
		$('table#qa-challenges').bootstrapTable({
		    url: '/api/challenges',
			search: true,
		    columns: [
				{ field: 'Name', title: 'Name' },
				{ field: 'Points',title: 'Points' },
			],
			onClickRow: function(row, $element, field) {
				window.location.href = '/solve/' + row.Id;
			},
		});
	},
	initProfile: function() {
		$('table#history').bootstrapTable({
		    url: '/api/history',
			showRefresh: true,
		    columns: [
				{ field: 'Id', title: 'Id' },
				{ field: 'Challenge',title: 'Challenge' },
				{ field: 'State', title: 'State',
					formatter: function(data, row, index) {
					    switch (data) {
					        case 1:
					            return '<span class="label label-success">Success</span>';
					            break
					        case 2:
					            return '<span class="label label-danger">Failed</span>';
					            break
							case 0:
								return '<span class="label label-info">In progress</span>';
								break
					        default:
					            return '<span class="label label-warning">Undefined</span>';
					    }
					},
				},
				{ field: 'Duration',title: 'Duration' },
			],
		});
	},
}