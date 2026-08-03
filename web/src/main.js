// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import './assets/iconfont/iconfont.css'
import './assets/css/tailwind.css'
import App from './App.vue'
import {createPinia} from "pinia"
import ElementPlus from "element-plus"
import {router} from "./router"
import 'bootstrap/dist/css/bootstrap.min.css'
import "element-plus/dist/index.css"
import 'vant/lib/index.css'


// Import Bootstrap JS
import 'bootstrap'

const app = createApp(App)
app.use(createPinia())
app.use(ElementPlus)
app.use(router).mount('#app')


