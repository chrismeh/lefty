let app = new Vue({
    el: '#app',
    data: {
        "products": [],
        "order": "price",
        "search": "",
    },
    methods: {
        fetchProducts: async function () {
            let url = `/api/products?order=${this.order}`
            let search = this.search.trim()
            if (search.length > 2) {
                url += `&search=${search}`
            }

            const response = await fetch(url);
            const json = await response.json();
            this.products = json.data;
        },
        updateSearchTerm: debounce(async function() {
            await this.fetchProducts()
        }, 500),
        resetSearchTerm: async function() {
            this.search = "";
            await this.fetchProducts();
        }
    },
    watch: {
        order: async function() {
            await this.fetchProducts();
        }
    },
    mounted: async function () {
       await this.fetchProducts();
    },
})

// Blatantly stolen from https://davidwalsh.name/javascript-debounce-function
function debounce(func, wait, immediate) {
    let timeout;
    return function() {
        let context = this, args = arguments;
        let later = function() {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        let callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
}