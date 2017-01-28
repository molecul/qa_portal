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
}