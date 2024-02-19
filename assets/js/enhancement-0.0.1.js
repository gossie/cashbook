document.addEventListener('DOMContentLoaded', () => {
    const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);

    // Add a click event on each of them
    $navbarBurgers.forEach( el => {
        el.addEventListener('click', () => {

            // Get the target from the "data-target" attribute
            const target = el.dataset.target;
            const $target = document.getElementById(target);

            // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
            el.classList.toggle('is-active');
            $target.classList.toggle('is-active');

        });
    });

    for (let i=0; i<document.forms.length; i++) {
        const f = document.forms.item(i);
        const requiredFields = f.querySelectorAll('input.validate-required');
        const submitButton = f.querySelector('button[type="submit"]');
        submitButton.addEventListener('click', (ev) => requiredFields.forEach(requiredField => validateRequired(requiredField, ev)));
    }
});


function validateRequired(requiredField, ev) {
    if (requiredField.value === '') {
        ev.preventDefault();
        if (!requiredField.classList.contains('is-danger')) {
            requiredField.classList.add('is-danger');
        }
    } else {
        requiredField.classList.remove('is-danger');
    }
}
