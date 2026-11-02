<template>
    <div
        class="admin-editor"
        :class="{ 'admin-editor--dark': isDark, }"
    >
        <!-- Toolbar -->
        <div class="admin-editor__toolbar">
            <button
                v-for="btn in toolbarButtons"
                :key="btn.title"
                type="button"
                class="admin-editor__toolbar-btn"
                :class="{ 'admin-editor__toolbar-btn--active': btn.active, }"
                :title="btn.title"
                @click="btn.action"
                v-html="btn.icon"
            />
            <span class="admin-editor__separator" />
            <button
                type="button"
                class="admin-editor__toolbar-btn"
                title="HTML"
                :class="{ 'admin-editor__toolbar-btn--active': showHTML, }"
                @click="toggleHTML"
            >
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="16"
                    height="16"
                >
                    <polyline points="16 18 22 12 16 6" />
                    <polyline points="8 6 2 12 8 18" />
                </svg>
            </button>
        </div>

        <!-- Editor area -->
        <div
            v-show="!showHTML"
            class="admin-editor__content"
        >
            <div
                ref="editorRef"
                class="admin-editor__editable"
                :contenteditable="!readonly"
                :placeholder="placeholder"
                @input="onInput"
                @paste="onPaste"
            />
        </div>

        <!-- HTML source -->
        <textarea
            v-show="showHTML"
            :value="modelValue"
            class="admin-editor__source"
            :placeholder="placeholder"
            @input="onHTMLInput"
        />

        <!-- Character count -->
        <div
            v-if="maxLength"
            class="admin-editor__footer"
        >
            <span
                class="admin-editor__count"
                :class="{ 'admin-editor__count--over': charCount > maxLength, }"
            >
                {{ charCount }} / {{ maxLength }}
            </span>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, computed, watch, onMounted, nextTick, } from 'vue';

    const props = withDefaults(
        defineProps<{
            modelValue: string
            placeholder?: string
            readonly?: boolean
            maxLength?: number
            minHeight?: string
            isDark?: boolean
        }>(),
        {
            modelValue: '',
            placeholder: 'Введите текст...',
            readonly: false,
            maxLength: 0,
            minHeight: '300px',
            isDark: false,
        }
    );

    const emit = defineEmits<{
        'update:modelValue': [value: string,]
    }>();

    const editorRef = ref<HTMLDivElement | null>(null,);
    const showHTML = ref(false,);

    const charCount = computed(() => props.modelValue.length,);

    const toolbarButtons = computed(() => [
        {
            title: 'Жирный',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M15.6 10.79c.97-.67 1.65-1.77 1.65-2.79 0-2.26-1.75-4-4-4H7v14h7.04c2.09 0 3.71-1.7 3.71-3.79 0-1.52-.86-2.82-2.15-3.42zM10 6.5h3c.83 0 1.5.67 1.5 1.5s-.67 1.5-1.5 1.5h-3v-3zm3.5 9H10v-3h3.5c.83 0 1.5.67 1.5 1.5s-.67 1.5-1.5 1.5z"/></svg>',
            action: () => execCmd('bold',),
            active: false,
        },
        {
            title: 'Курсив',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M10 4v3h2.21l-3.42 8H6v3h8v-3h-2.21l3.42-8H18V4z"/></svg>',
            action: () => execCmd('italic',),
            active: false,
        },
        {
            title: 'Подчёркнутый',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M12 17c3.31 0 6-2.69 6-6V3h-2.5v8c0 1.93-1.57 3.5-3.5 3.5S8.5 12.93 8.5 11V3H6v8c0 3.31 2.69 6 6 6zm-7 2v2h14v-2H5z"/></svg>',
            action: () => execCmd('underline',),
            active: false,
        },
        {
            title: 'Зачёркнутый',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M10 19h4v-3h-4v3zM5 4v3h5v3h4V7h5V4H5zM3 14h18v-2H3v2z"/></svg>',
            action: () => execCmd('strikeThrough',),
            active: false,
        },
        {
            title: 'Заголовок 1',
            icon: '<span style="font-weight:700;font-size:14px;">H1</span>',
            action: () => execCmd('formatBlock', '<h1>',),
            active: false,
        },
        {
            title: 'Заголовок 2',
            icon: '<span style="font-weight:700;font-size:13px;">H2</span>',
            action: () => execCmd('formatBlock', '<h2>',),
            active: false,
        },
        {
            title: 'Заголовок 3',
            icon: '<span style="font-weight:700;font-size:12px;">H3</span>',
            action: () => execCmd('formatBlock', '<h3>',),
            active: false,
        },
        {
            title: 'Маркированный список',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M4 10.5c-.83 0-1.5.67-1.5 1.5s.67 1.5 1.5 1.5 1.5-.67 1.5-1.5-.67-1.5-1.5-1.5zm0-6c-.83 0-1.5.67-1.5 1.5s.67 1.5 1.5 1.5 1.5-.67 1.5-1.5-.67-1.5-1.5-1.5zm0 12c-.83 0-1.5.67-1.5 1.5s.67 1.5 1.5 1.5 1.5-.67 1.5-1.5-.67-1.5-1.5-1.5zM7 19h14v-2H7v2zm0-6h14v-2H7v2zm0-8v2h14V5H7z"/></svg>',
            action: () => execCmd('insertUnorderedList',),
            active: false,
        },
        {
            title: 'Нумерованный список',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M2 17h2v.5H3v1h1v.5H2v1h3v-4H2v1zm1-9h1V4H2v1h1v3zm-1 3h1.8L2 13.1v.9h3v-1H3.2L5 10.9V10H2v1zm5-6v2h14V5H7zm0 14h14v-2H7v2zm0-6h14v-2H7v2z"/></svg>',
            action: () => execCmd('insertOrderedList',),
            active: false,
        },
        {
            title: 'Цитата',
            icon:
                '<svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16"><path d="M6 17h3l2-4V7H5v6h3zm8 0h3l2-4V7h-6v6h3z"/></svg>',
            action: () => execCmd('formatBlock', '<blockquote>',),
            active: false,
        },
        {
            title: 'Ссылка',
            icon:
                '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/></svg>',
            action: () => insertLink(),
            active: false,
        },
    ]);

    // Sync modelValue to editor content
    watch(
        () => props.modelValue,
        (newVal,) => {
            if (editorRef.value && editorRef.value.innerHTML !== newVal && !showHTML.value) {
                editorRef.value.innerHTML = newVal;
            }
        }
    );

    onMounted(() => {
        if (editorRef.value) {
            editorRef.value.innerHTML = props.modelValue;
        }
    });

    function execCmd(command: string, value?: string,) {
        document.execCommand(command, false, value,);
        editorRef.value?.focus();
        updateModel();
    }

    function insertLink() {
        const url = window.prompt('Введите URL ссылки:',);
        if (url) {
            execCmd('createLink', url,);
        }
    }

    function onInput() {
        updateModel();
    }

    function onPaste(e: ClipboardEvent,) {
        e.preventDefault();
        const text = e.clipboardData?.getData('text/plain',);
        if (text) {
            document.execCommand('insertText', false, text,);
        }
    }

    function onHTMLInput(e: Event,) {
        const val = (e.target as HTMLTextAreaElement).value;
        emit('update:modelValue', val,);
    }

    function toggleHTML() {
        showHTML.value = !showHTML.value;
        if (!showHTML.value && editorRef.value) {
            nextTick(() => {
                if (editorRef.value) {
                    editorRef.value.innerHTML = props.modelValue;
                }
            });
        }
    }

    function updateModel() {
        if (editorRef.value) {
            emit('update:modelValue', editorRef.value.innerHTML,);
        }
    }
</script>

<style scoped>
.admin-editor {
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 12px;
    overflow: hidden;
    background: var(--admin-surface, #ffffff);
    transition: border-color 0.2s, background 0.3s;
}

.admin-editor--dark {
    background: var(--admin-surface, #1a1d29);
    border-color: var(--admin-border, #2d2d3d);
}

.admin-editor:focus-within {
    border-color: var(--admin-primary, #6c5ce7);
    box-shadow: 0 0 0 3px rgba(108, 92, 231, 0.1);
}

/* Toolbar */
.admin-editor__toolbar {
    display: flex;
    align-items: center;
    gap: 2px;
    padding: 8px 12px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    flex-wrap: wrap;
    background: var(--admin-bg, #f0f2f5);
}

.admin-editor--dark .admin-editor__toolbar {
    background: rgba(255, 255, 255, 0.03);
    border-color: var(--admin-border, #2d2d3d);
}

.admin-editor__toolbar-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    border-radius: 6px;
    background: transparent;
    color: var(--admin-text-secondary, #636e72);
    cursor: pointer;
    transition: all 0.15s;
}

.admin-editor__toolbar-btn:hover {
    background: rgba(0, 0, 0, 0.05);
    color: var(--admin-text-primary, #2d3436);
}

.admin-editor--dark .admin-editor__toolbar-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-editor__toolbar-btn--active {
    background: rgba(108, 92, 231, 0.12);
    color: var(--admin-primary, #6c5ce7);
}

.admin-editor__separator {
    width: 1px;
    height: 24px;
    background: var(--admin-border, #e0e0e0);
    margin: 0 4px;
}

.admin-editor--dark .admin-editor__separator {
    background: var(--admin-border, #2d2d3d);
}

/* Editable content */
.admin-editor__content {
    padding: 0;
}

.admin-editor__editable {
    min-height: v-bind(minHeight);
    padding: 16px;
    outline: none;
    font-size: 14px;
    line-height: 1.7;
    color: var(--admin-text-primary, #2d3436);
    overflow-y: auto;
}

.admin-editor--dark .admin-editor__editable {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-editor__editable:empty::before {
    content: attr(placeholder);
    color: var(--admin-text-secondary, #636e72);
    pointer-events: none;
}

.admin-editor__editable :deep(h1) {
    font-size: 24px;
    font-weight: 700;
    margin: 16px 0 8px;
}

.admin-editor__editable :deep(h2) {
    font-size: 20px;
    font-weight: 700;
    margin: 14px 0 6px;
}

.admin-editor__editable :deep(h3) {
    font-size: 17px;
    font-weight: 600;
    margin: 12px 0 4px;
}

.admin-editor__editable :deep(blockquote) {
    border-left: 3px solid var(--admin-primary, #6c5ce7);
    padding: 8px 16px;
    margin: 12px 0;
    background: rgba(108, 92, 231, 0.05);
    border-radius: 0 8px 8px 0;
    color: var(--admin-text-secondary, #636e72);
}

.admin-editor__editable :deep(ul),
.admin-editor__editable :deep(ol) {
    padding-left: 24px;
    margin: 8px 0;
}

.admin-editor__editable :deep(a) {
    color: var(--admin-primary, #6c5ce7);
    text-decoration: underline;
}

.admin-editor__editable :deep(img) {
    max-width: 100%;
    border-radius: 8px;
    margin: 12px 0;
}

/* HTML source */
.admin-editor__source {
    width: 100%;
    min-height: v-bind(minHeight);
    padding: 16px;
    border: none;
    outline: none;
    resize: vertical;
    font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
    line-height: 1.6;
    background: var(--admin-bg, #f0f2f5);
    color: var(--admin-text-primary, #2d3436);
    tab-size: 2;
}

.admin-editor--dark .admin-editor__source {
    background: rgba(255, 255, 255, 0.03);
    color: var(--admin-text-primary, #e0e0e0);
}

/* Footer */
.admin-editor__footer {
    display: flex;
    justify-content: flex-end;
    padding: 8px 16px;
    border-top: 1px solid var(--admin-border, #e0e0e0);
}

.admin-editor--dark .admin-editor__footer {
    border-color: var(--admin-border, #2d2d3d);
}

.admin-editor__count {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-editor__count--over {
    color: #e17055;
    font-weight: 600;
}
</style>
