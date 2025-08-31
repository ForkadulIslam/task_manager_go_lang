<template>
  <div class="rich-text-editor">
    <div v-if="editor" class="editor-toolbar">
      <!-- Formatting Buttons -->
      <button @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.196-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.783-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path></svg>
      </button>
      <button @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1m0-1c0-2.21 1.79-4 4-4h4m-4 4v1m0-1c0-2.21 1.79-4 4-4h4m-4 4v1m0-1c0-2.21 1.79-4 4-4h4m-4 4v1m0-1c0-2.21 1.79-4 4-4h4"></path></svg>
      </button>
      <button @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </button>

      <!-- Heading Buttons -->
      <button @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }">H1</button>
      <button @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }">H2</button>
      <button @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }">H3</button>

      <!-- List Buttons -->
      <button @click="editor.chain().focus().toggleBulletList().run()" :class="{ 'is-active': editor.isActive('bulletList') }">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
      </button>
      <button @click="editor.chain().focus().toggleOrderedList().run()" :class="{ 'is-active': editor.isActive('orderedList') }">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"></path></svg>
      </button>

      <!-- Other Buttons -->
      <button @click="editor.chain().focus().toggleBlockquote().run()" :class="{ 'is-active': editor.isActive('blockquote') }">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path></svg>
      </button>
      <button @click="editor.chain().focus().setHorizontalRule().run()">HR</button>
      <button @click.prevent="addImage">Add Image</button>
    </div>
    <editor-content :editor="editor" />
  </div>
</template>

<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import Image from '@tiptap/extension-image';
import ImageResize from 'tiptap-extension-resize-image';
import { defineEmits, defineProps, watch } from 'vue';
import apiClient from '../services/api';

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:modelValue']);

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    Image,
    ImageResize,
  ],
  editorProps: {
    handlePaste(view, event, slice) {
      const items = (event.clipboardData || event.originalEvent.clipboardData).items;
      for (const item of items) {
        if (item.type.indexOf('image') === 0) {
          const file = item.getAsFile();
          uploadImage(file);
          return true; // Prevent default paste behavior
        }
      }
      return false;
    },
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML());
  },
});

watch(() => props.modelValue, (value) => {
  const isSame = editor.value.getHTML() === value;
  if (isSame) {
    return;
  }
  editor.value.commands.setContent(value, false);
});

const uploadImage = async (file) => {
  if (file) {
    const formData = new FormData();
    formData.append('attachment', file);
    try {
      const response = await apiClient.post('/upload-attachment', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      const url = `${import.meta.env.VITE_API_BASE_URL}/${response.data.path}`;
      if (url) {
        editor.value.chain().setImage({ src: url }).run();
      }
    } catch (error) {
      console.error('Image upload failed:', error);
    }
  }
};

const addImage = async () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*';
  input.onchange = async (e) => {
    const file = e.target.files[0];
    uploadImage(file);
  };
  input.click();
};

</script>

<style>
.rich-text-editor .ProseMirror {
  border: 1px solid #4a5568;
  padding: 0.5rem;
  min-height: 200px;
  background-color: #1a202c;
  color: #e2e8f0;
  border-radius: 0.375rem;
}

.rich-text-editor .editor-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  border: 1px solid #4a5568;
  border-bottom: none;
  padding: 0.5rem;
  background-color: #2d3748;
  border-top-left-radius: 0.375rem;
  border-top-right-radius: 0.375rem;
}

.rich-text-editor .editor-toolbar button {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 0.25rem;
  margin-bottom: 0.25rem;
  padding: 0.5rem;
  border: 1px solid transparent;
  border-radius: 0.25rem;
  background-color: transparent;
  color: #e2e8f0;
  transition: background-color 0.2s;
}

.rich-text-editor .editor-toolbar button:hover {
  background-color: #4a5568;
}

.rich-text-editor .editor-toolbar button.is-active {
  background-color: #2b6cb0;
  color: white;
}

.rich-text-editor img {
  max-width: 100%;
  height: auto;
}

.rich-text-editor .ProseMirror-selectednode {
  outline: 3px solid #63b3ed;
}
</style>