export const useTableParser = () => {
    const parseHTMLTableToJSON = (htmlString) => {
      const parser = new DOMParser()
      const doc = parser.parseFromString(htmlString, 'text/html')
      const table = doc.querySelector('table')
      
      if (!table) return { columns: [], rows: [] }
      
      const allRows = Array.from(table.querySelectorAll('tr'))
      if (allRows.length === 0) return { columns: [], rows: [] }
      
      // Track which cells are occupied by rowspan/colspan
      const occupiedCells = new Map() // Map of "rowIdx,colIdx" -> true
      
      // Find maximum number of columns
      let maxColumns = 0
      
      allRows.forEach((row, rowIdx) => {
        const cells = row.querySelectorAll('td, th')
        let colIdx = 0
        
        cells.forEach(cell => {
          // Skip occupied cells
          while (occupiedCells.has(`${rowIdx},${colIdx}`)) {
            colIdx++
          }
          
          const colspan = parseInt(cell.getAttribute('colspan')) || 1
          const rowspan = parseInt(cell.getAttribute('rowspan')) || 1
          
          // Mark cells as occupied by this cell's rowspan/colspan
          for (let r = 0; r < rowspan; r++) {
            for (let c = 0; c < colspan; c++) {
              occupiedCells.set(`${rowIdx + r},${colIdx + c}`, true)
            }
          }
          
          colIdx += colspan
          maxColumns = Math.max(maxColumns, colIdx)
        })
      })
      
      console.log(`Max columns detected: ${maxColumns}`)
      
      // Create column definitions
      const columns = Array.from({ length: maxColumns }, (_, idx) => ({
        title: `Column ${idx + 1}`,
        field: `col${idx}`,
        editor: 'input',
        width: 150,
      }))
      
      // Extract ALL rows with proper colspan/rowspan handling
      const rows = []
      const cellOccupancy = new Map() // Reset for data extraction
      
      allRows.forEach((tr, rowIdx) => {
        const row = {}
        const cells = tr.querySelectorAll('td, th')
        let colIdx = 0
        
        cells.forEach(cell => {
          // Skip columns occupied by previous rowspan cells
          while (cellOccupancy.has(`${rowIdx},${colIdx}`)) {
            colIdx++
          }
          
          const colspan = parseInt(cell.getAttribute('colspan')) || 1
          const rowspan = parseInt(cell.getAttribute('rowspan')) || 1
          const text = cell.textContent.trim()
          
          // Mark cells occupied by this cell's rowspan/colspan
          for (let r = 0; r < rowspan; r++) {
            for (let c = 0; c < colspan; c++) {
              cellOccupancy.set(`${rowIdx + r},${colIdx + c}`, true)
              
              // Only put text in the first cell of the span
              if (r === 0 && c === 0 && colIdx + c < maxColumns) {
                row[`col${colIdx + c}`] = text
              } else if (colIdx + c < maxColumns) {
                row[`col${colIdx + c}`] = '' // Empty for spanned cells
              }
            }
          }
          
          colIdx += colspan
        })
        
        rows.push(row)
      })
      
      console.log('Parsed table:', { 
        headerCount: columns.length, 
        rowCount: rows.length,
        sampleRow: rows[0],
        lastRow: rows[rows.length - 1]
      })
      
      return { columns, rows }
    }
    
    const convertJSONToHTML = (columns, rows) => {
      let html = '<table><tbody>'
      
      rows.forEach(row => {
        html += '<tr>'
        columns.forEach(col => {
          html += `<td>${row[col.field] || ''}</td>`
        })
        html += '</tr>'
      })
      html += '</tbody></table>'
      
      return html
    }
    
    return { parseHTMLTableToJSON, convertJSONToHTML }
  }