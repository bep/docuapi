export class Highlight {
	constructor(
		opts = {
			contentSelector: '.content',
			markClass: 'da-highlight-mark',
		}
	) {
		this.opts = opts;
		this.nodeStack = [];
	}

	apply(re) {
		const treeWalker = document.createTreeWalker(
			this.content(),
			NodeFilter.SHOW_TEXT,
			{
				acceptNode: (node) => {
					if (node.parentNode.classList.contains(this.opts.markClass)) {
						return NodeFilter.FILTER_REJECT;
					}
					if (/\S/.test(node.data)) {
						return NodeFilter.FILTER_ACCEPT;
					}
					return NodeFilter.FILTER_REJECT;
				},
			},
			false
		);

		// Firs pass: Collect all the text nodes matching the provided regexp.
		// TODO(bep) improve text matching.
		let matches = [];
		while (treeWalker.nextNode()) {
			let node = treeWalker.currentNode;
			if (node.data.match(re)) {
				matches.push(node);
			}
		}

		// Second pass: Replace the matches with nodes wrapped in <mark> tags.
		matches.forEach((node) => {
			// Clone the parent so we can restore it.
			let parentClone = node.parentNode.cloneNode(true);

			parentClone.childNodes.forEach((node) => {
				if (node.nodeType !== Node.TEXT_NODE) {
					return;
				}

				let match = node.data.match(re);
				if (!match) {
					return;
				}

				let pos = node.data.indexOf(match[0], match.index);
				if (pos === -1) {
					return;
				}

				let mark = document.createElement('mark');
				mark.classList.add(this.opts.markClass);

				let wordNode = node.splitText(pos);
				wordNode.splitText(match[0].length);
				let wordClone = wordNode.cloneNode(true);

				mark.appendChild(wordClone);
				parentClone.replaceChild(mark, wordNode);
			});

			if (node.parentNode && node.parentNode.parentNode) {
				this.nodeStack.push({
					old: node.parentNode,
					new: parentClone,
				});
				node.parentNode.parentNode.replaceChild(parentClone, node.parentNode);
			}
		});
	}

	remove() {
		while (this.nodeStack.length) {
			let pair = this.nodeStack.pop();
			pair.new.parentNode.replaceChild(pair.old, pair.new);
		}
	}

	content() {
		return document.querySelector(this.opts.contentSelector);
	}
}
