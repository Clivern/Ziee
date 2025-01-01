<template>
  <div class="space-y-4">
    <p v-if="isNewVersion" class="text-sm text-theme-textLight rounded-md border border-theme-border bg-theme-hover px-3 py-2">
      {{ $t('prompts_page.immutable_note') }}
    </p>

    <div v-if="showName">
      <label :for="`${idPrefix}-name`" class="form-label">{{ $t('prompts_page.name') }}</label>
      <input
        :id="`${idPrefix}-name`"
        :value="form.name"
        type="text"
        required
        class="input-field"
        :placeholder="$t('prompts_page.name_placeholder')"
        @input="updateField('name', $event.target.value)"
      />
    </div>

    <div>
      <label :for="`${idPrefix}-type`" class="form-label">{{ $t('prompts_page.type') }}</label>
      <select
        :id="`${idPrefix}-type`"
        :value="form.type"
        class="input-field"
        @change="updateField('type', $event.target.value)"
      >
        <option value="text">{{ $t('prompts_page.type_text') }}</option>
        <option value="chat">{{ $t('prompts_page.type_chat') }}</option>
      </select>
    </div>

    <div>
      <label :for="`${idPrefix}-content`" class="form-label">{{ $t('prompts_page.content') }}</label>
      <p class="text-xs text-theme-textLight mb-1.5">
        {{ $t('prompts_page.variables_hint_prefix') }}
        <span class="font-mono text-theme-text">{{ PROMPT_VARIABLE_EXAMPLE }}</span>.
        {{ $t('prompts_page.variables_hint_suffix') }}
      </p>
      <textarea
        :id="`${idPrefix}-content`"
        :value="form.content"
        required
        rows="8"
        class="input-field font-mono text-sm"
        :placeholder="`${$t('prompts_page.content_placeholder')}\n${PROMPT_VARIABLE_EXAMPLE}`"
        @input="updateField('content', $event.target.value)"
      />
      <p v-if="variables.length" class="mt-2 text-xs text-theme-textLight">
        {{ $t('prompts_page.variables_available') }}
        <span v-for="v in variables" :key="v" class="ml-1 font-mono text-theme-text">{{ v }}</span>
      </p>
    </div>

    <div>
      <label :for="`${idPrefix}-config`" class="form-label">{{ $t('prompts_page.config') }}</label>
      <p class="text-xs text-theme-textLight mb-1.5">{{ $t('prompts_page.config_hint') }}</p>
      <textarea
        :id="`${idPrefix}-config`"
        :value="form.config"
        rows="4"
        class="input-field font-mono text-sm"
        placeholder="{}"
        @input="updateField('config', $event.target.value)"
      />
    </div>

    <label class="flex items-center gap-2 text-sm text-theme-text cursor-pointer">
      <input
        :checked="form.production"
        type="checkbox"
        class="rounded border-theme-border text-primary-600 focus:ring-primary-800"
        @change="updateField('production', $event.target.checked)"
      />
      {{ $t('prompts_page.set_production') }}
    </label>

    <div>
      <label :for="`${idPrefix}-commit`" class="form-label">{{ $t('prompts_page.commit_message') }}</label>
      <textarea
        :id="`${idPrefix}-commit`"
        :value="form.commitMessage"
        rows="2"
        class="input-field"
        :placeholder="$t('prompts_page.commit_message_placeholder')"
        @input="updateField('commitMessage', $event.target.value)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { extractPromptVariables, emptyPromptForm, PROMPT_VARIABLE_EXAMPLE } from '@/lib/prompt'

const form = defineModel({
  type: Object,
  default: emptyPromptForm,
})

defineProps({
  showName: {
    type: Boolean,
    default: false,
  },
  isNewVersion: {
    type: Boolean,
    default: false,
  },
  idPrefix: {
    type: String,
    default: 'prompt',
  },
})

function updateField(key, value) {
  form.value = { ...form.value, [key]: value }
}

const variables = computed(() => extractPromptVariables(form.value?.content))
</script>
