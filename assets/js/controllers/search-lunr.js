import * as lunr from 'js/lib/lunr.js';
import { Highlight } from 'js/helpers';

function nextUntil(elem, selector) {
	var siblings = [];

	elem = elem.nextElementSibling;

	while (elem) {
		if (elem.matches(selector)) break;
		siblings.push(elem);
		elem = elem.nextElementSibling;
	}

	return siblings;
}

export function newSearchController() {
	var index;

	var buildIndex = function (config) {
		var builder = new lunr.Builder();

		builder.pipeline.add(lunr.trimmer, lunr.stopWordFilter, lunr.stemmer);

		builder.searchPipeline.add(lunr.stemmer);

		config.call(builder, builder);
		return builder.build();
	};

	function populateIndex() {
		index = buildIndex(function () {
			this.ref('id');
			this.field('title', { boost: 10 });
			this.field('body');
			this.pipeline.add(lunr.trimmer, lunr.stopWordFilter);

			document.querySelectorAll('h1, h2').forEach((headerEl) => {
				let body = '';
				nextUntil(headerEl, 'h1, h2').forEach((el) => {
					body = body.concat(' ', el.textContent);
				});
				this.add({
					id: headerEl.id,
					title: headerEl.textContent,
					body: body,
				});
			});
		});
	}

	let highlight = new Highlight();

	return {
		query: '',
		results: [],
		init: function () {
			return this.$nextTick(() => {
				populateIndex();
				this.$watch('query', () => {
					this.search();
				});
			});
		},
		search: function () {
			highlight.remove();
			let results = index.search(this.query).filter((item) => item.score > 0.0001);

			this.results = results.map((item) => {
				var elem = document.getElementById(item.ref);

				return {
					title: elem.innerText,
					id: item.ref,
				};
			});

			if (this.results.length > 0) {
				highlight.apply(new RegExp(this.query, 'i'));
			}
		},
	};
}
