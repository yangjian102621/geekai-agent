<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">基础设置</h5>
    </div>
    <div class="card-body">
      <div class="container pt-2">
        <form class="form">
          <div class="mb-3">
            <label class="form-label">网站标题</label>
            <input type="text" class="form-control" v-model="config.title" />
          </div>

          <div class="mb-3">
            <label class="form-label">控制台标题</label>
            <input
              type="text"
              class="form-control"
              v-model="config.admin_title"
            />
          </div>

          <div class="mb-3">
            <label class="form-label">网站Slogan</label>
            <input type="text" class="form-control" v-model="config.slogan" />
          </div>

          <div class="mb-3">
            <label class="form-label"
              >网站Logo<el-tooltip
                effect="dark"
                content="未授权的系统不允许修改 Logo"
                raw-content
                placement="right"
              >
                <i class="iconfont icon-info ms-1"></i> </el-tooltip
            ></label>
            <div class="d-flex gap-2">
              <input type="text" class="form-control" v-model="config.logo" />
              <el-upload
                :auto-upload="true"
                :show-file-list="false"
                :http-request="uploadLogo"
              >
                <button type="button" class="btn btn-primary">
                  <i class="iconfont icon-upload"></i>
                </button>
              </el-upload>
            </div>
          </div>

          <div class="mb-3">
            <label class="form-label"
              >网站版权<el-tooltip
                effect="dark"
                content="未授权的系统不允许修改版权信息"
                raw-content
                placement="right"
              >
                <i class="iconfont icon-info ms-1"></i>
              </el-tooltip>
            </label>
            <div class="d-flex gap-2">
              <input
                type="text"
                class="form-control"
                v-model="config.copyright"
              />
            </div>
          </div>

          <div class="mb-3">
            <label class="form-label"
              >开放注册
              <el-tooltip
                effect="dark"
                content="关闭注册之后只能通过管理后台添加用户"
                raw-content
                placement="right"
              >
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <div>
              <el-switch size="large" v-model="config.enabled_register" />
            </div>
          </div>

          <div class="mb-3">
            <label class="form-label">邮件域名白名单</label>
            <div class="input-group mb-3">
              <items-input :value="config.email_white_list" />
            </div>
          </div>

          <div class="mb-3">
            <label class="form-label">客服微信二维码</label>
            <div class="d-flex gap-2">
              <input
                type="text"
                class="form-control"
                v-model="config.wechat_card_url"
              />
              <el-upload
                :auto-upload="true"
                :show-file-list="false"
                :http-request="uploadWechatCard"
              >
                <button type="button" class="btn btn-primary">
                  <i class="iconfont icon-upload"></i>
                </button>
              </el-upload>
            </div>
          </div>

          <div class="mb-3">
            <label class="form-label"
              >默认聊天应用
              <el-tooltip
                effect="dark"
                content="用户聊天页面默认应用，如果为空则使用第一个应用"
                raw-content
                placement="right"
              >
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-select
              v-model="config.app_id"
              class="w-100"
              placeholder="请选择应用"
              filterable
              clearable
            >
              <el-option
                v-for="item in appList"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <div class="d-flex align-items-center">
                  <el-image
                    :src="item.icon"
                    style="height: 24px; width: 24px"
                    class="rounded-circle me-2"
                  >
                    <template #error>
                      <i class="iconfont icon-image"></i>
                    </template>
                  </el-image>
                  <span>{{ item.name }}</span>
                  <el-tag class="ms-auto text-muted text-sm"
                    >{{ item.score }}积分/次</el-tag
                  >
                </div>
              </el-option>
            </el-select>
          </div>

          <div class="mb-3">
            <label class="form-label">注册赠送积分</label>
            <input
              type="number"
              class="form-control"
              v-model.number="config.init_score"
            />
          </div>

          <div class="mb-3">
            <label class="form-label">每日签到赠送积分</label>
            <input
              type="number"
              class="form-control"
              v-model.number="config.daily_score"
            />
          </div>

          <div class="mb-3">
            <label class="form-label">邀请新用户赠送积分</label>
            <input
              type="number"
              class="form-control"
              v-model.number="config.invite_score"
            />
          </div>

          <div class="d-flex gap-2 justify-content-center pt-3">
            <button type="button" class="btn btn-primary" @click="saveConfig">
              提交
            </button>
            <button
              type="button"
              class="btn btn-light"
              @click="config = copyObj(configBackup)"
            >
              重置
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import ItemsInput from '@/components/ItemsInput.vue'
  import { useAdminSystemStore } from '@/js/store/admin/system'
  import { copyObj } from '@/js/utils/libs.js'

  const systemStore = useAdminSystemStore()
  const { config, configBackup, appList } = storeToRefs(systemStore)
  const {
    saveConfig,
    uploadLogo,
    uploadWechatCard,
    uploadBotAvatar,
    uploadUserAvatar,
    initialize,
  } = systemStore

  onMounted(() => {
    initialize()
  })
</script>

<style scoped lang="scss">
  .card {
    padding: 0.5rem 2rem 1rem 2rem;
  }
</style>
