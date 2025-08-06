<template>
    <div class="bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-200 font-sans">
        <main class="container mx-auto px-6 md:px-16 py-12 md:py-12">
            <div class="lg:flex lg:gap-6 lg:justify-center">

                <aside class="lg:w-1/4 xl:w-1/5 lg:pr-8">
                    <div class="sticky top-24">
                        <h3 class="text-sm font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 mb-4">On this page</h3>
                        <nav class="space-y-1">
                            <a v-for="link in navLinks" :key="link.href" :href="link.href" @click.prevent="scrollToSection(link.href)"
                                :class="[
                                    'block px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200',
                                    activeSection === link.href.substring(1)
                                        ? 'bg-blue-50 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300'
                                        : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
                                ]">
                                {{ link.text }}
                            </a>
                        </nav>
                    </div>
                </aside>

                <article class="prose-lg dark:prose-invert max-w-none mt-12 lg:mt-0 lg:w-3/4 xl:w-4/5">

                    <section id="getting-started" ref="sections" class="scroll-mt-24 pb-16 border-b border-gray-200 dark:border-gray-800">
                        <h2 class="tracking-tight text-2xl font-bold !mb-2">Getting Started</h2>
                        <p class="!mb-3">
                            Welcome to the Go-API Mart documentation! <br/> This guide will walk you through the essential steps to start using the APIs available on our marketplace. All APIs on our platform share a consistent structure for authentication and requests.
                        </p>
                        <ol class="list-decimal list-inside flex flex-col   gap-2">
                            <li><strong>Create an Account:</strong> If you haven't already, <router-link to="/register">sign up for a free account</router-link>.</li>
                            <li><strong>Find an API:</strong> <router-link to="/browse">Browse our catalog</router-link> to find an API that suits your needs.</li>
                            <li><strong>Subscribe:</strong> On the API's detail page, click the "Subscribe" button. This will grant you access to the API.</li>
                            <li><strong>Get Your API Key:</strong> After subscribing, you can find your unique API key in your <router-link to="/dashboard">consumer dashboard</router-link> under the "My Subscriptions" section.</li>
                        </ol>
                    </section>

                    <section id="authentication" ref="sections" class="scroll-mt-24 py-16 border-b border-gray-200 dark:border-gray-800">
                        <h2 class="tracking-tight text-2xl font-bold !mb-2">Authentication</h2>
                        <p class="!mb-6">
                            All requests to any API on the Go-API Mart platform must be authenticated using an Api-key scheme. You need to include your API key in the <code>Authorization</code> header for every request. You can find your personal API key for each subscription in your dashboard.
                        </p>
                        <pre class="!bg-gray-50 dark:!bg-gray-800/50 !rounded-xl !mb-6"><code class="language-http hljs">Api-key: YOUR_API_KEY</code></pre>
                        
                        <div class="mt-8 p-4 rounded-lg bg-amber-50 dark:bg-amber-900/30 border border-amber-200 dark:border-amber-500/30 flex items-start gap-4 not-prose">
                                <svg class="h-6 w-6 text-amber-500 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126z" />
                                </svg>
                                <div class="text-sm text-amber-800 dark:text-amber-200">
                                        <strong class="font-semibold">Warning:</strong> Keep your API keys secure and do not expose them in client-side code. It is recommended to use them from your backend server.
                                </div>
                        </div>
                    </section>

                    <section id="making-requests" ref="sections" class="scroll-mt-24 py-16 border-b border-gray-200 dark:border-gray-800">
                        <h2 class="tracking-tight text-2xl font-bold !mb-2">Making Requests</h2>
                        <p class="!mb-6">
                            The base URL for every API is dynamic and constructed based on the provider's username. The structure is always:
                        </p>
                        <pre class="!bg-gray-50 dark:!bg-gray-800/50 !rounded-xl"><code class="language-text hljs">https://{provider_username}.zxsttm.tech</code></pre>
                        <p class="!mt-6 !mb-6">
                            You can find the specific base URL and available endpoints on the detail page for each API.
                        </p>
                        <h4 class="!mt-10">Example Request with cURL</h4>
                        <p class="!mb-6">
                            Here is an example of how to make a request to the Geolocation API, assuming the provider's username is <code>seller_store</code>.
                        </p>
                        <pre class="!bg-gray-50 dark:!bg-gray-800/50 !rounded-xl"><code class="language-bash hljs">curl --request GET \
         --url 'https://seller_store.api.zxsttm.tech/geolocation/v1/ip/8.8.8.8' \
         `--header 'Api-key: YOUR_API_KEY'`</code></pre>
                    </section>

                    <section id="error-handling" ref="sections" class="scroll-mt-24 py-16">
                        <h2 class="tracking-tight text-2xl !mb-2">Error Handling</h2>
                        <p class="!mb-6">
                            Our API gateway uses standard HTTP status codes to indicate the success or failure of a request.
                        </p>
                        <div class="not-prose my-8 overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700">
                                <table class="w-full text-left text-sm">
                                        <thead class="bg-gray-50 dark:bg-gray-800 text-xs font-semibold uppercase tracking-wider text-gray-600 dark:text-gray-300">
                                                <tr>
                                                        <th class="px-6 py-3">Status Code</th>
                                                        <th class="px-6 py-3">Meaning</th>
                                                        <th class="px-6 py-3">Description</th>
                                                </tr>
                                        </thead>
                                        <tbody class="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-900">
                                                <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>200 OK</code></td>
                                                        <td class="px-6 py-4">Success</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">The request was successful.</td>
                                                </tr>
                                                <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>400 Bad Request</code></td>
                                                        <td class="px-6 py-4">Client Error</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">The request was malformed, such as missing a required parameter.</td>
                                                </tr>
                                                <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>401 Unauthorized</code></td>
                                                        <td class="px-6 py-4">Authentication Error</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">Your API key is missing, invalid, or expired.</td>
                                                </tr>
                                                <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>403 Forbidden</code></td>
                                                        <td class="px-6 py-4">Permission Denied</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">You are not subscribed to this API or your plan does not allow this action.</td>
                                                </tr>
                                                <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>404 Not Found</code></td>
                                                        <td class="px-6 py-4">Endpoint Not Found</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">The requested endpoint does not exist.</td>
                                                </tr>
                                                 <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>429 Too Many Requests</code></td>
                                                        <td class="px-6 py-4">Rate Limit Exceeded</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">You have exceeded the allowed number of requests. Please wait and try again.</td>
                                                </tr>
                                                <tr>
                                                        <td class="px-6 py-4 font-mono text-gray-800 dark:text-gray-200"><code>500 Internal Server Error</code></td>
                                                        <td class="px-6 py-4">Server Error</td>
                                                        <td class="px-6 py-4 text-gray-600 dark:text-gray-400">Something went wrong on our end. Please try again later.</td>
                                                </tr>
                                        </tbody>
                                </table>
                        </div>
                    </section>

                </article>
            </div>
        </main>
    </div>
</template>

<script setup>
import { onMounted, onUnmounted, nextTick, ref } from 'vue';
import { RouterLink } from 'vue-router';

const navLinks = ref([
        { href: '#getting-started', text: 'Getting Started' },
        { href: '#authentication', text: 'Authentication' },
        { href: '#making-requests', text: 'Making Requests' },
        { href: '#error-handling', text: 'Error Handling' },
]);

const sections = ref([]);
const activeSection = ref('');
let observer;

onMounted(() => {
        // Set the first section as active by default
        if (navLinks.value.length > 0) {
                activeSection.value = navLinks.value[0].href.substring(1);
        }

        nextTick(() => {
                if(window.feather) window.feather.replace();
                if(window.hljs) window.hljs.highlightAll();

                // Observer for active section highlighting
                const observerOptions = {
                        root: null,
                        rootMargin: '0px 0px -50% 0px', // Trigger when section is in the top half of the viewport
                        threshold: 0.5,
                };

                observer = new IntersectionObserver((entries) => {
                        entries.forEach((entry) => {
                                if (entry.isIntersecting) {
                                        activeSection.value = entry.target.id;
                                }
                        });
                }, observerOptions);

                const sectionElements = document.querySelectorAll('section');
                sectionElements.forEach((section) => {
                        observer.observe(section);
                });
        });
});

onUnmounted(() => {
        if (observer) {
                observer.disconnect();
        }
});

const scrollToSection = (hash) => {
        const element = document.querySelector(hash);
        if (element) {
                element.scrollIntoView({ behavior: 'smooth' });
                // Manually set active section on click for instant feedback
                activeSection.value = hash.substring(1); 
        }
};
</script>

<style scoped>
/* Scoped styles can be added here if needed, but Tailwind handles most of it. */
strong {
    font-weight: 600;
}
</style>