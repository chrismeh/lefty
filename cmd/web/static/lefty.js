let app = new Vue({
    el: '#app',
    data: {
        "products": [],
        "pagination": null,
        "order": "price",
        "search": "",
        "retailer": "",
        "requestedPage": 1,
    },
    methods: {
        fetchProducts: async function () {
            const response = await fetch(this.buildApiUrl());
            const json = await response.json();
            this.products = json.data;
            this.pagination = json.meta;
        },
        resetSearchTerm: async function() {
            this.search = "";
            await this.fetchProducts();
        },
        buildApiUrl: function() {
            let url = `/api/products?order=${this.order}`
            let search = this.search.trim()
            if (search.length > 2) {
                url += `&search=${search}`
            }

            if (this.retailer !== "") {
                url += `&retailer=${this.retailer}`
            }

            if (this.requestedPage > 1) {
                url += `&page=${this.requestedPage}`
            }
            return url
        }
    },
    watch: {
        order: async function() {
            await this.fetchProducts();
        },
        retailer: async function() {
            await this.fetchProducts();
        },
        search: debounce(async function() {
            await this.fetchProducts();
        }, 500),
        requestedPage: async function() {
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