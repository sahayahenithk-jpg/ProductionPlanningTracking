<template>
  <div class="page">
    <section class="card">
      <div class="dashboard-shell">
        <aside class="dashboard-side">
          <div class="dashboard-side-header">
            <h3>Dashboard</h3>
            <p class="subtitle">Switch panels and manage your account.</p>
          </div>

          <div class="tabs-row">
            <button
              v-for="tab in tabs"
              :key="tab"
              :class="['tab-button', activeTab === tab ? 'active' : '']"
              @click="activeTab = tab"
            >
              {{ tab }}
            </button>
          </div>

          <div class="side-panel">
            <div class="profile-card compact side-profile-card">
              <p class="profile-heading">Profile</p>
              <p><strong>{{ user?.name || '-' }}</strong></p>
              <p>{{ user?.email || '-' }}</p>
            </div>
            <button class="secondary logout side-logout" @click="logout">Logout</button>
          </div>
        </aside>

        <div class="dashboard-main">
          <div class="header-row">
            <div>
              <h2>Dashboard</h2>
              <p class="subtitle">Performance overview, charts, and quick reports for your production workflow.</p>
            </div>
          </div>

          <div v-if="error" class="error">{{ error }}</div>

          <div v-if="activeTab === 'Overview'">
            <div class="cards-grid">
              <div class="summary-card">
                <div class="card-label">Total Products</div>
                <div class="card-value">{{ overview.totalProducts }}</div>
              </div>
              <div class="summary-card">
                <div class="card-label">Total Plans</div>
                <div class="card-value">{{ overview.totalPlans }}</div>
              </div>
              <div class="summary-card">
                <div class="card-label">Today's Production</div>
                <div class="card-value">{{ today.productionTotal }}</div>
              </div>
              <div class="summary-card">
                <div class="card-label">Achievement %</div>
                <div class="card-value">{{ today.achievementPct.toFixed(2) }}%</div>
              </div>
            </div>

            <div class="overview-panels">
              <div class="panel-card">
                <h3>Performance Score</h3>
                <div class="metric-grid">
                  <div class="metric-card">
                    <span class="metric-label">Produced</span>
                    <strong>{{ metrics.producedTotal }}</strong>
                  </div>
                  <div class="metric-card">
                    <span class="metric-label">Planned</span>
                    <strong>{{ metrics.plannedTotal }}</strong>
                  </div>
                  <div class="metric-card">
                    <span class="metric-label">Avg Daily</span>
                    <strong>{{ metrics.avgDaily }}</strong>
                  </div>
                  <div class="metric-card">
                    <span class="metric-label">Avg Achievement</span>
                    <strong>{{ metrics.avgAchievement.toFixed(1) }}%</strong>
                  </div>
                </div>
              </div>

              <div class="panel-card actions-panel">
                <h3>Quick actions</h3>
                <div class="overview-actions">
                  <router-link to="/products" class="action-button">Products</router-link>
                  <router-link to="/plans" class="action-button">Plans</router-link>
                  <router-link to="/variance" class="action-button">Reports</router-link>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'Performance'">
            <div class="filter-card">
              <div class="filter-title">
                <h3>Chart filters</h3>
                <span>Use filters to focus the dashboard data.</span>
              </div>
              <form @submit.prevent="applyFilters" class="filter-form">
                <div class="filter-grid">
                  <select v-model="filters.productId">
                    <option value="">All Products</option>
                    <option
                      v-for="product in products"
                      :key="product.productId"
                      :value="product.productId"
                    >
                      {{ product.productName }}
                    </option>
                  </select>
                  <input type="date" v-model="filters.startDate" placeholder="Start date" />
                  <input type="date" v-model="filters.endDate" placeholder="End date" />
                  <input type="month" v-model="filters.month" />
                </div>

                <div class="filter-actions">
                  <button type="submit" class="primary">Apply</button>
                  <button type="button" class="secondary" @click="resetFilters">Reset</button>
                </div>
              </form>
            </div>

            <div class="cards-grid">
              <div class="summary-card">
                <div class="card-label">Total Products</div>
                <div class="card-value">{{ overview.totalProducts }}</div>
              </div>
              <div class="summary-card">
                <div class="card-label">Total Plans</div>
                <div class="card-value">{{ overview.totalPlans }}</div>
              </div>
              <div class="summary-card">
                <div class="card-label">Today's Production</div>
                <div class="card-value">{{ today.productionTotal }}</div>
              </div>
              <div class="summary-card">
                <div class="card-label">Achievement %</div>
                <div class="card-value">{{ today.achievementPct.toFixed(2) }}%</div>
              </div>
            </div>

            <div class="charts-grid">
              <div class="chart-card">
                <div class="chart-header">
                  <h3>Production Trend</h3>
                  <span>{{ chartRangeLabel }}</span>
                </div>
                <div v-if="productionTrend.length" class="chart-list">
                  <div v-for="item in productionTrend" :key="item.label" class="chart-row">
                    <div class="chart-row-title">{{ item.label }}</div>
                    <div class="chart-bar">
                      <div class="chart-bar-fill" :style="{ width: item.width + '%' }"></div>
                    </div>
                    <div class="chart-row-value">{{ item.value }}</div>
                  </div>
                </div>
                <div v-else class="empty-state">No trend data available for the selected filters.</div>
              </div>

              <div class="chart-card">
                <div class="chart-header">
                  <h3>Product Achievement</h3>
                  <span>Average achievement by product</span>
                </div>
                <div v-if="productPerformance.length" class="chart-list">
                  <div v-for="item in productPerformance" :key="item.name" class="chart-row">
                    <div class="chart-row-title">{{ item.name }}</div>
                    <div class="chart-bar">
                      <div class="chart-bar-fill achievement" :style="{ width: item.width + '%' }"></div>
                    </div>
                    <div class="chart-row-value">{{ item.achievementPct.toFixed(1) }}%</div>
                  </div>
                </div>
                <div v-else class="empty-state">Apply filters or add production data to populate this chart.</div>
              </div>

              <div class="chart-card">
                <div class="chart-header">
                  <h3>Production Distribution</h3>
                  <span>Share of output by product</span>
                </div>
                <div v-if="distribution.length" class="chart-list">
                  <div v-for="item in distribution" :key="item.name" class="chart-row">
                    <div class="chart-row-title">{{ item.name }}</div>
                    <div class="chart-bar">
                      <div class="chart-bar-fill distribution" :style="{ width: item.width + '%' }"></div>
                    </div>
                    <div class="chart-row-value">{{ item.percent.toFixed(0) }}%</div>
                  </div>
                </div>
                <div v-else class="empty-state">No distribution data available yet.</div>
              </div>

              <div class="chart-card summary-metrics-card">
                <div class="chart-header">
                  <h3>Production Score</h3>
                  <span>Key filtered performance metrics</span>
                </div>
                <div class="metric-grid">
                  <div class="metric-card">
                    <span class="metric-label">Produced</span>
                    <strong>{{ metrics.producedTotal }}</strong>
                  </div>
                  <div class="metric-card">
                    <span class="metric-label">Planned</span>
                    <strong>{{ metrics.plannedTotal }}</strong>
                  </div>
                  <div class="metric-card">
                    <span class="metric-label">Avg Daily</span>
                    <strong>{{ metrics.avgDaily }}</strong>
                  </div>
                  <div class="metric-card">
                    <span class="metric-label">Avg Achievement</span>
                    <strong>{{ metrics.avgAchievement.toFixed(1) }}%</strong>
                  </div>
                </div>
              </div>

              <div class="chart-card">
                <div class="chart-header">
                  <h3>Top Plan Compliance</h3>
                  <span>Compare actual vs planned totals</span>
                </div>
                <div v-if="topCompliance.length" class="chart-list">
                  <div v-for="item in topCompliance" :key="item.name" class="chart-row compliance-row">
                    <div class="chart-row-title">{{ item.name }}</div>
                    <div class="compliance-values">
                      <div><strong>{{ item.produced }}</strong><span>Produced</span></div>
                      <div><strong>{{ item.planned }}</strong><span>Planned</span></div>
                    </div>
                  </div>
                </div>
                <div v-else class="empty-state">No compliance results to display.</div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'Reports'" class="report-layout">
            <div class="report-card">
              <div class="report-card-header">
                <h3>Top Plan Compliance</h3>
                <span>Actual vs planned totals</span>
              </div>
              <div v-if="topCompliance.length" class="chart-list">
                <div v-for="item in topCompliance" :key="item.name" class="chart-row compliance-row">
                  <div class="chart-row-title">{{ item.name }}</div>
                  <div class="compliance-values">
                    <div><strong>{{ item.produced }}</strong><span>Produced</span></div>
                    <div><strong>{{ item.planned }}</strong><span>Planned</span></div>
                  </div>
                </div>
              </div>
              <div v-else class="empty-state">No compliance results to display.</div>
            </div>

            <div class="report-card">
              <div class="report-card-header">
                <h3>Quick actions</h3>
                <span>Move faster through your workflows</span>
              </div>
              <div class="actions-row report-actions">
                <router-link v-if="role === 'admin'" to="/products" class="action-button">Manage Products</router-link>
                <router-link v-if="role !== 'operator'" to="/plans" class="action-button">Create Plans</router-link>
                <router-link v-if="role !== 'operator'" to="/variance" class="action-button">View Report</router-link>
                <router-link v-if="role === 'operator'" to="/production" class="action-button">Enter Production</router-link>
              </div>
              <div class="profile-card compact report-profile">
                <p><strong>Name:</strong> {{ user?.name || '-' }}</p>
                <p><strong>Email:</strong> {{ user?.email || '-' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'
import { getUserRole } from '../services/auth'

const router = useRouter()
const role = ref(getUserRole())
const user = ref(null)
const error = ref('')
const overview = ref({ totalProducts: 0, totalPlans: 0 })
const today = ref({ productionTotal: 0, achievementPct: 0 })
const products = ref([])
const productionTrend = ref([])
const productPerformance = ref([])
const distribution = ref([])
const topCompliance = ref([])
const metrics = ref({ producedTotal: 0, plannedTotal: 0, avgDaily: 0, avgAchievement: 0 })
const tabs = computed(() => {
  if (role.value === 'operator') {
    return ['Overview', 'Performance']
  }
  return ['Overview', 'Performance', 'Reports']
})
const activeTab = ref('Overview')

const formatToday = () => new Date().toISOString().split('T')[0]
const getLastWeek = () => {
  const date = new Date()
  date.setDate(date.getDate() - 6)
  return date.toISOString().split('T')[0]
}

const filters = ref({ productId: '', startDate: getLastWeek(), endDate: formatToday(), month: '' })

const chartRangeLabel = computed(() => {
  if (filters.value.month) {
    return filters.value.month
  }
  if (filters.value.startDate && filters.value.endDate) {
    return `${filters.value.startDate} → ${filters.value.endDate}`
  }
  return 'Last 7 days'
})

const loadProfile = async () => {
  try {
    const response = await api.get('/user')
    user.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load profile'
  }
}

const loadProducts = async () => {
  try {
    const response = await api.get('/products')
    products.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load products'
  }
}

const loadOverview = async () => {
  try {
    const response = await api.get('/reports/summary')
    overview.value = {
      totalProducts: response.data.totalProducts,
      totalPlans: response.data.totalPlans,
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load overview'
  }
}

const loadTodayStats = async () => {
  try {
    const response = await api.get('/reports/summary', {
      params: {
        startDate: formatToday(),
        endDate: formatToday(),
      },
    })
    today.value = {
      productionTotal: response.data.productionTotal,
      achievementPct: response.data.achievementPct,
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load today summary'
  }
}

const buildParams = () => {
  const params = {}
  if (filters.value.productId) params.productId = filters.value.productId
  if (filters.value.month) params.month = filters.value.month
  if (filters.value.startDate) params.startDate = filters.value.startDate
  if (filters.value.endDate) params.endDate = filters.value.endDate
  return params
}

const processVarianceData = (rows) => {
  const byDate = {}
  const byProduct = {}
  let producedTotal = 0
  let plannedTotal = 0
  let totalAchievement = 0

  rows.forEach((row) => {
    producedTotal += row.producedQuantity
    plannedTotal += row.plannedQuantity
    totalAchievement += row.achievementPct

    byDate[row.productionDate] = (byDate[row.productionDate] || 0) + row.producedQuantity

    if (!byProduct[row.productName]) {
      byProduct[row.productName] = { produced: 0, planned: 0, achievementSum: 0, count: 0 }
    }
    byProduct[row.productName].produced += row.producedQuantity
    byProduct[row.productName].planned += row.plannedQuantity
    byProduct[row.productName].achievementSum += row.achievementPct
    byProduct[row.productName].count += 1
  })

  const dateKeys = Object.keys(byDate).sort()
  const maxTrend = Math.max(...Object.values(byDate), 1)
  productionTrend.value = dateKeys.map((label) => ({
    label,
    value: byDate[label],
    width: Math.round((byDate[label] / maxTrend) * 100),
  }))

  const productEntries = Object.entries(byProduct)
    .map(([name, data]) => ({
      name,
      produced: data.produced,
      planned: data.planned,
      achievementPct: data.count ? data.achievementSum / data.count : 0,
    }))
    .sort((a, b) => b.achievementPct - a.achievementPct)

  const maxAchievement = Math.max(...productEntries.map((item) => item.achievementPct), 100)
  productPerformance.value = productEntries.slice(0, 5).map((item) => ({
    ...item,
    width: Math.round((item.achievementPct / maxAchievement) * 100),
  }))

  const totalProducedForShare = Math.max(producedTotal, 1)
  distribution.value = productEntries
    .slice(0, 5)
    .map((item) => ({
      name: item.name,
      value: item.produced,
      percent: (item.produced / totalProducedForShare) * 100,
      width: Math.round((item.produced / totalProducedForShare) * 100),
    }))

  topCompliance.value = productEntries
    .slice(0, 4)
    .map((item) => ({
      name: item.name,
      produced: item.produced,
      planned: item.planned,
    }))

  const uniqueDays = dateKeys.length || 1
  metrics.value = {
    producedTotal,
    plannedTotal,
    avgDaily: Math.round(producedTotal / uniqueDays),
    avgAchievement: rows.length ? totalAchievement / rows.length : 0,
  }
}

const loadCharts = async () => {
  error.value = ''
  try {
    const params = buildParams()
    const response = await api.get('/reports/variance', { params })
    processVarianceData(response.data)
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load chart data'
  }
}

const applyFilters = async () => {
  await loadCharts()
}

const resetFilters = async () => {
  filters.value = { productId: '', startDate: getLastWeek(), endDate: formatToday(), month: '' }
  await loadCharts()
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/')
}

onMounted(async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/')
    return
  }
  await Promise.all([loadProfile(), loadProducts(), loadOverview(), loadTodayStats(), loadCharts()])
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 30px;
  display: flex;
  justify-content: center;
  background: linear-gradient(135deg, #eef2f3, #d9e4f5);
  font-family: Arial, sans-serif;
}

.card {
  width: 100%;
  max-width: 1100px;
  background: white;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.12);
}

.dashboard-shell {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 24px;
}

.dashboard-side {
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 24px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.dashboard-side-header h3 {
  margin: 0 0 10px;
  font-size: 20px;
  color: #111827;
}

.side-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.side-profile-card {
  background: white;
  border: 1px solid #dbeafe;
  padding: 18px;
}

.profile-heading {
  margin: 0 0 10px;
  color: #475569;
  font-size: 13px;
}

.side-logout {
  width: 100%;
  text-align: center;
}

.dashboard-main {
  min-width: 0;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.header-row h2 {
  font-size: 26px;
  margin: 0;
  color: #111827;
}

.subtitle {
  color: #6b7280;
  margin: 4px 0 0;
}

.tabs-row {
  display: grid;
  gap: 12px;
}

.tab-button {
  width: 100%;
  text-align: left;
  padding: 14px 16px;
  border: 1px solid #d1d5db;
  border-radius: 14px;
  background: white;
  color: #334155;
  font-weight: 600;
  cursor: pointer;
}

.tab-button.active {
  background: #1d4ed8;
  color: white;
  border-color: #1d4ed8;
}

.filter-card {
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  background: #f8fafc;
  padding: 20px;
  margin-bottom: 24px;
}

.filter-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 18px;
}

.filter-title h3 {
  margin: 0;
  color: #111827;
}

.filter-title span {
  color: #4b5563;
}

.filter-form {
  display: grid;
  gap: 16px;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(180px, 1fr));
  gap: 14px;
}

.filter-grid select,
.filter-grid input {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #d1d5db;
  border-radius: 12px;
  background: white;
  font-size: 14px;
}

.filter-actions {
  display: flex;
  gap: 12px;
}

.primary,
.secondary,
.logout {
  border: none;
  border-radius: 10px;
  padding: 12px 18px;
  cursor: pointer;
}

.primary {
  background: #1d4ed8;
  color: white;
}

.secondary {
  background: #e2e8f0;
  color: #0f172a;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card {
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
}

.card-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
}

.card-value {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
}

.overview-panels {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 18px;
}

.panel-card,
.report-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  padding: 24px;
}

.actions-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.overview-actions,
.report-actions {
  display: grid;
  gap: 12px;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(320px, 1fr));
  gap: 18px;
  margin-bottom: 24px;
}

.chart-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  padding: 22px;
  min-height: 280px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 18px;
}

.chart-header h3 {
  margin: 0;
  font-size: 18px;
  color: #111827;
}

.chart-header span {
  color: #6b7280;
  font-size: 13px;
}

.chart-list {
  display: grid;
  gap: 12px;
}

.chart-row {
  display: grid;
  grid-template-columns: 1.4fr 3fr 0.9fr;
  gap: 12px;
  align-items: center;
}

.chart-row-title {
  font-size: 14px;
  color: #374151;
}

.chart-bar {
  height: 14px;
  background: #e5e7eb;
  border-radius: 999px;
  overflow: hidden;
}

.chart-bar-fill {
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, #2563eb, #60a5fa);
}

.chart-bar-fill.achievement {
  background: linear-gradient(90deg, #16a34a, #6ee7b7);
}

.chart-bar-fill.distribution {
  background: linear-gradient(90deg, #f97316, #fb923c);
}

.chart-row-value {
  text-align: right;
  font-size: 14px;
  font-weight: 700;
  color: #111827;
}

.empty-state {
  color: #6b7280;
  padding: 20px 0;
  font-size: 14px;
}

.summary-metrics-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(140px, 1fr));
  gap: 16px;
}

.metric-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 16px;
}

.metric-label {
  display: block;
  color: #6b7280;
  font-size: 13px;
  margin-bottom: 8px;
}

.metric-card strong {
  font-size: 24px;
  color: #111827;
}

.compliance-row {
  grid-template-columns: 1.4fr 2.5fr;
}

.compliance-values {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
}

.compliance-values div {
  text-align: right;
}

.compliance-values span {
  display: block;
  color: #6b7280;
  font-size: 12px;
}

.report-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 18px;
}

.report-card-header h3 {
  margin: 0;
  font-size: 18px;
  color: #111827;
}

.report-card-header span {
  color: #475569;
  font-size: 13px;
}

.section {
  margin-top: 20px;
}

.profile-card {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 18px;
}

.profile-card p {
  margin: 8px 0;
  color: #111827;
}

.actions-row {
  margin-top: 28px;
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
}

.action-button {
  background: #2563eb;
  color: white;
  padding: 12px 18px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
}

.logout {
  background: #ef4444;
  color: white;
  cursor: pointer;
}

.error {
  color: #b91c1c;
  margin-bottom: 16px;
}

@media (max-width: 980px) {
  .dashboard-shell {
    grid-template-columns: 1fr;
  }

  .overview-panels,
  .charts-grid,
  .cards-grid,
  .filter-grid,
  .overview-actions,
  .report-actions {
    grid-template-columns: 1fr;
  }
}
</style>
