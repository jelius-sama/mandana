(function() {
    function getThemePreference() {
        return localStorage.getItem('themePreference') || 'system';
    }

    function applyTheme() {
        const preference = getThemePreference();
        let isDark = false;

        if (preference === 'system') {
            isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        } else {
            isDark = preference === 'dark';
        }

        document.documentElement.classList.toggle('dark', isDark);
        document.dispatchEvent(new CustomEvent('theme-changed'));
    }

    function cycleTheme() {
        const current = getThemePreference();
        let next;

        if (current === 'system') {
            const isDarkNow = window.matchMedia('(prefers-color-scheme: dark)').matches;
            next = isDarkNow ? 'light' : 'dark';
        } else {
            next = current === 'light' ? 'dark' : 'light';
        }

        localStorage.setItem('themePreference', next);
        applyTheme();
    }

    applyTheme();

    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
        if (getThemePreference() === 'system') {
            applyTheme();
        }
    });

    document.addEventListener('click', (e) => {
        const themeSwitcher = e.target.closest('[data-theme-switcher]');
        if (themeSwitcher) {
            e.preventDefault();
            cycleTheme();
        }
    });
})();
