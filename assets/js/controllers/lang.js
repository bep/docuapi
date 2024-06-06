const debug = 0 ? console.log.bind(console, '[lang]') : function () {};

const toggleCodeblockVisibility = function (lang, visible) {
	debug('toggleCodeblockVisibility', lang, visible);
	document.querySelectorAll(`.highlight code.language-${lang}`).forEach((el) => {
		let highlight = el.closest('.highlight');
		highlight.style.display = visible ? 'block' : 'none';
	});
};

export function newLangController() {
	return {
		tabs: [],
        contentChildren: [],
		changeLanguage: function (index) {
			debug('changeLanguage', index);
            const anchor = this.contentChildren.find(el => {
                return el.getBoundingClientRect().top > 0
            })
            const anchorValue = anchor.getBoundingClientRect().top;
			for (let i = 0; i < this.tabs.length; i++) {
				let isActive = i === index;
				this.tabs[i].active = isActive;

				toggleCodeblockVisibility(this.tabs[i].key, isActive);
			}
            const updatedAnchorElementValue = anchor.getBoundingClientRect().top;
            anchorValue !== updatedAnchorElementValue && window.scrollBy(0, updatedAnchorElementValue);
		},
		initLangs: function (tabs) {
			debug('initLangs', tabs);
			tabs[0].active = true;
			this.tabs = tabs;
            this.contentChildren = Array.from(document.querySelector('.page-wrapper .content').children)
                .filter(el => el.tagName !== 'BLOCKQUOTE')
                .filter(el => !el.classList.contains('highlight'))

			return this.$nextTick(() => {
				this.changeLanguage(0);
			});
		},
	};
}
