<template>
  <div class="container-fluid p-4">
    <ParamConfig
      v-model="mockParams"
      title="参数配置测试"
      @file-upload="handleFileUpload"
      @preview="handlePreview"
    />

    <el-image-viewer
      v-if="previewURL"
      :url-list="[previewURL]"
      @close="previewURL = ''"
    />
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import ParamConfig from '@/components/admin/ParamConfig.vue'
  import { showMessageOK } from '@/js/utils/dialog.js'
  import { adminUploadFile } from '@/js/store/common.js'
  import { showLoading, closeLoading } from '@/js/utils/dialog.js'

  const previewURL = ref('')

  const mockParams = ref([
    {
      label: '客户姓名',
      name: 'customer_name',
      type: 'String',
      default: '张三',
      placeholder: '请输入客户姓名',
      required: true,
    },
    {
      label: '客户简介',
      name: 'customer_brief',
      type: 'String',
      default: '这是一位重点客户，需要定期跟进。',
      placeholder: '请输入客户简介',
      required: false,
    },
    {
      label: '客户年龄',
      name: 'customer_age',
      type: 'Number',
      default: 28,
      required: false,
    },
    {
      label: '是否会员',
      name: 'is_vip',
      type: 'Boolean',
      default: true,
      required: true,
    },
    {
      label: '回访时间',
      name: 'next_visit_at',
      type: 'DateTime',
      default: '2025-01-08 10:00:00',
      required: false,
    },
    {
      label: '偏好渠道',
      name: 'preferred_channel',
      type: 'Select',
      default: '微信',
      options: ['微信', '电话', '邮件'],
      placeholder: '请选择偏好渠道',
      required: true,
    },
    {
      label: '客户星级',
      name: 'customer_star',
      type: 'Radio',
      default: '五星',
      options: ['一星', '三星', '五星'],
      required: true,
    },
    {
      label: '可选服务',
      name: 'optional_services',
      type: 'CheckBox',
      default: ['咨询', '培训'],
      options: ['咨询', '培训', '实施'],
      required: false,
    },
    {
      label: '合同图片',
      name: 'contract_image',
      type: 'Image',
      default: 'https://dummyimage.com/240x240/409eff/ffffff&text=Image',
      max_filesize: 5,
      required: false,
    },
    {
      label: '欢迎语音',
      name: 'welcome_audio',
      type: 'Audio',
      default: 'https://www.w3schools.com/html/horse.mp3',
      max_filesize: 10,
      required: false,
    },
    {
      label: '演示视频',
      name: 'demo_video',
      type: 'Video',
      default: 'https://www.w3schools.com/html/mov_bbb.mp4',
      max_filesize: 30,
      required: false,
    },
    {
      label: '附件文档',
      name: 'attachment_doc',
      type: 'Doc',
      default:
        'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf',
      max_filesize: 20,
      required: false,
    },
  ])

  const handleFileUpload = (file, param) => {
    showLoading('正在上传...')
    adminUploadFile(file, (data) => {
      param.default = data.url
      closeLoading()
      showMessageOK('上传成功！')
    })
  }

  const handlePreview = (url) => {
    previewURL.value = url
  }
</script>

<style lang="scss" scoped>
  // Base class
  .callout {
    padding: 1.25rem;
    margin-top: 1.25rem;
    margin-bottom: 1.25rem;
    color: var(--bd-callout-color, inherit);
    background-color: var(--bd-callout-bg, var(--bs-gray-100));
    border-left: 0.25rem solid var(--bd-callout-border, var(--bs-gray-300));
  }

  // Modifier classes
  .callout-info {
    --bd-callout-color: #087990;
    --bd-callout-bg: var(--bs-info-bg-subtle);
    --bd-callout-border: #9eeaf9;
  }

  .callout-warning {
    --bd-callout-color: var(--bs-warning-text);
    --bd-callout-bg: var(--bs-warning-bg-subtle);
    --bd-callout-border: var(--bs-warning-border-subtle);
  }

  .callout-danger {
    --bd-callout-color: var(--bs-danger-text);
    --bd-callout-bg: var(--bs-danger-bg-subtle);
    --bd-callout-border: var(--bs-danger-border-subtle);
  }
</style>
