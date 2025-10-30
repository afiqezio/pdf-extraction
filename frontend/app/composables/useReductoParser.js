export const useReductoParser = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || 'http://localhost:8080'

  const parseReductoJSON = async (reductoJSON) => {
    try {
      const response = await fetch(`${apiBase}/api/v1/parse-reducto`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(reductoJSON),
      })

      if (!response.ok) {
        const error = await response.json()
        throw new Error(error.error || 'Failed to parse Reducto JSON')
      }

      const data = await response.json()
      console.log('Parsed Reducto data:', data)
      return data
    } catch (error) {
      console.error('Error parsing Reducto JSON:', error)
      throw error
    }
  }

  return {
    parseReductoJSON,
  }
}
