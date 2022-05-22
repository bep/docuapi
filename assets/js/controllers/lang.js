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
		changeLanguage: function (index) {
			debug('changeLanguage', index);
			for (let i = 0; i < this.tabs.length; i++) {
				let isActive = i === index;
				this.tabs[i].active = isActive;

				toggleCodeblockVisibility(this.tabs[i].key, isActive);
			}
		},
		initLangs: function (tabs) {
			debug('initLangs', tabs);
			tabs[0].active = true;
			this.tabs = tabs;

			return this.$nextTick(() => {
				this.changeLanguage(0);
			});
		},
	};
}
