function switchStyle() {
    var currentStyle = document.getElementById('style_switcher').getAttribute('href');
    var button = document.getElementById('themeToggleButton');
    if (currentStyle === '/static/dark_theme.css') {
        document.getElementById('style_switcher').setAttribute('href', '/static/light_theme.css');
        localStorage.setItem('theme', 'light');
        button.textContent = 'Dark mode';
    } else {
        document.getElementById('style_switcher').setAttribute('href', '/static/dark_theme.css');
        localStorage.setItem('theme', 'dark');
        button.textContent = 'Light mode';
    }
}

function loadTheme() {
    var theme = localStorage.getItem('theme');
    var button = document.getElementById('themeToggleButton');
    if (theme === 'light') {
        document.getElementById('style_switcher').setAttribute('href', '/static/light_theme.css');
        button.textContent = 'Dark mode';
    } else {
        document.getElementById('style_switcher').setAttribute('href', '/static/dark_theme.css');
        button.textContent = 'Light mode';
    }
}

window.onload = loadTheme;