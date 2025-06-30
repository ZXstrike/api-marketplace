<template>
  <div v-if="endpoint" class="flex flex-col md:flex-row gap-x-8 gap-y-12">
    <div class="w-full space-y-6">
      <h1 class="text-3xl font-bold font-mono">Description</h1>
      <div v-if="renderedDocumentation" class="max-w-none py-3" v-html="renderedDocumentation"></div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';

const props = defineProps({
  endpoint: { type: Object, required: true },
  baseUrl: { type: String, required: true }
});

const docString = computed(() => {
  let curl = `curl --request ${props.endpoint.http_method}
`;
  curl += `--url '${props.baseUrl}${props.endpoint.path}'
`;
  curl += `--header 'Token: YOUR_API_KEY'`;
  
  const requestExample = `## Base Request Example\n\n\`\`\`bash\n${curl}\n\`\`\``;

  return `${requestExample}\n\n${props.endpoint.documentation || ''}`;
});

const marked = new Marked(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    }
  })
);

const renderedDocumentation = computed(() => {
  if (docString.value) {
    console.log('Rendering documentation:', docString.value);
    return marked.parse(docString.value);
  }
  return '';
});

</script>

<style scoped>
:deep(pre) {
  background-color: #282c34; /* Corresponds to atom-one-dark background */
  color: #abb2bf; /* Corresponds to atom-one-dark foreground */
  padding: 1em;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1rem;
}

:deep(p){
  margin: 0.5rem 0;
  line-height: 1.6rem;
}

:deep(code) {
  font-family: 'Courier New', Courier, monospace;
}

:deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5em 0;
  border-radius: 0.5rem;
  overflow: hidden;
}

:deep(thead) {
  background-color: #1f2937;
  color: #f9fafb;
  font-weight: bold;
}

:deep(th),
:deep(td) {
  padding: 0.75rem 1rem;
  border: 1px solid #374151;
  text-align: left;
}

:deep(tbody tr) {
  background-color: #374151;
  color: #d1d5db;
}

:deep(tbody tr:nth-child(even)) {
    background-color: #4b5563;
}

:deep(td code) {
    background-color: #1f2937;
    padding: 0.2em 0.4em;
    margin: 0;
    font-size: 85%;
    border-radius: 3px;
}
</style>