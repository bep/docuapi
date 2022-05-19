'use strict';

const debug = 0 ? console.log.bind(console, '[toc]') : function () {};

const headerEls = () => document.querySelectorAll('.content h1, .content h2, .content h3');

const setProgress = function (self, el) {
	let mainEl = document.querySelector('.content');
	let mainHeight = mainEl.offsetHeight;
	let mainStart = mainEl.offsetTop;
	let progress = Math.round(((el.offsetTop - mainStart) / mainHeight) * 100);
	self.activeHeading.title = el.innerText;
	self.activeHeading.progress = progress;
};

export function newToCController() {
	const setOpenRecursive = function (row, shouldOpen) {
		if (!row.sub) {
			return false;
		}
		let isOpen = false;

		row.sub.forEach((rowsub) => {
			rowsub.open = shouldOpen(rowsub);
			rowsub.active = rowsub.open;

			rowsub.open = rowsub.open || setOpenRecursive(rowsub, shouldOpen);

			if (rowsub.open) {
				isOpen = true;
			}
		});

		row.active_parent = isOpen;
		isOpen = isOpen || shouldOpen(row);
		row.open = isOpen;
		row.active = shouldOpen(row);

		return isOpen;
	};

	return {
		activeHeading: {
			title: '',
			progress: 0,
		},
		showHeading: true,
		rows: [],
		load: function (rows) {
			this.rows = rows;
		},

		transitions: function () {
			return {
				'x-transition:enter.duration.500ms': '',
				'x-transition:leave.duration.400ms': '',
				'x-transition.scale.origin.top.left.80': '',
			};
		},

		rowClass: function (row) {
			return {
				['x-bind:class']() {
					return `toc-h${row.level}${row.active ? ' active' : ''}${row.active_parent ? ' active-parent' : ''}`;
				},
			};
		},

		click: function (row) {
			this.rows.forEach((row2) => {
				setOpenRecursive(row2, (row3) => {
					return row === row3;
				});
			});
		},

		onScroll: function () {
			debug('onScroll');
			let scrollpos = window.scrollY;

			headerEls().forEach((el) => {
				let offset = el.offsetTop;

				if (offset < scrollpos + 10) {
					debug('Set for', el.id);

					this.rows.forEach((row) => {
						setOpenRecursive(row, (row) => {
							return row.id === el.id;
						});
					});

					setProgress(this, el);
				}

				if (window.innerHeight + scrollpos >= document.body.offsetHeight) {
					this.activeHeading.progress = 100;
				}
			});
		},
	};
}
