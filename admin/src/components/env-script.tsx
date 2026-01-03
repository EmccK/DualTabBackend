// 运行时环境变量注入脚本
// 标准做法: 通过 window.__ENV__ 注入运行时配置

export function EnvScript() {
  // 服务端读取环境变量
  const env = {
    API_URL: process.env.API_URL || '',
  }

  return (
    <script
      dangerouslySetInnerHTML={{
        __html: `window.__ENV__ = ${JSON.stringify(env)}`,
      }}
    />
  )
}
