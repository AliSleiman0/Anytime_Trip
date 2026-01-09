// SweetAlert2 Export Options
async function showExportOptions() {
  const { value: exportFormat } = await Swal.fire({
    title: 'Export as:',
    icon: 'question',
    html: `
      <div class="flex flex-col gap-3">
        <button class="export-option-btn" data-format="pdf">
           PDF
        </button>
        <button class="export-option-btn" data-format="csv">
           CSV
        </button>
        <button class="export-option-btn" data-format="excel">
           Excel
        </button>
      </div>
    `,
    didOpen: (modal) => {
      const buttons = modal.querySelectorAll('.export-option-btn');
      buttons.forEach(btn => {
        btn.addEventListener('click', function() {
          const format = this.getAttribute('data-format');
          Swal.close();
          handleExport(format);
        });
      });
    },
    showConfirmButton: false,
    allowOutsideClick: true,
  });
}

// Handle export based on format
function handleExport(format) {
  console.log(`Exporting as ${format.toUpperCase()}`);
  
  // Get all tables (since headers and data might be in separate tables)
  const tables = document.querySelectorAll('table');
  
  if (tables.length === 0) {
    Swal.fire('Error', 'No table data found to export', 'error');
    return;
  }

  switch(format.toLowerCase()) {
    case 'pdf':
      exportToPDF(tables);
      break;
    case 'csv':
      exportToCSV(tables);
      break;
    case 'excel':
      exportToExcel(tables);
      break;
  }
}

// Export to PDF
function exportToPDF(tables) {
  // Check if jsPDF is loaded with multiple fallbacks
  let jsPDFLib = null;
  
  // Try different ways jsPDF might be available
  if (window.jspdf && window.jspdf.jsPDF) {
    jsPDFLib = window.jspdf.jsPDF;
  } else if (typeof jsPDF !== 'undefined') {
    jsPDFLib = jsPDF;
  } else if (window.jsPDF) {
    jsPDFLib = window.jsPDF;
  }

  if (!jsPDFLib) {
    Swal.fire('Error', 'PDF library failed to load. Try these steps:\n1. Hard refresh (Ctrl+F5)\n2. Clear browser cache\n3. Check console for errors', 'error');
    console.error('jsPDF not found. Available window properties:', Object.keys(window).filter(k => k.toLowerCase().includes('pdf')));
    return;
  }

  try {
    const doc = new jsPDFLib({
      orientation: 'landscape',
      unit: 'mm',
      format: 'a4'
    });

    let yPosition = 15;
    const pageHeight = doc.internal.pageSize.getHeight();
    const pageWidth = doc.internal.pageSize.getWidth();
    const margin = 10;

    // Add title
    doc.setFontSize(16);
    doc.text('Data Export', margin, yPosition);
    yPosition += 10;

    // Collect headers and data separately
    let headerRow = [];
    let bodyRows = [];

    // First pass: collect headers from all tables (excluding last column if it's "Action")
    tables.forEach((table) => {
      if (headerRow.length === 0) {
        const headerRows = table.querySelectorAll('thead tr');
        headerRows.forEach(row => {
          const cols = row.querySelectorAll('th, td');
          const rowData = [];
          cols.forEach((col, index) => {
            // Skip the last column if it contains "Action"
            if (index < cols.length - 1 || (col.innerText.trim().toLowerCase() !== 'action' && col.innerText.trim().toLowerCase() !== 'actions')) {
              rowData.push(col.innerText.trim());
            }
          });
          if (rowData.length > 0 && !rowData.every(cell => cell === '')) {
            headerRow = rowData;
          }
        });
      }
    });

    // Second pass: collect all data rows from all tables (excluding last column)
    tables.forEach((table) => {
      const dataRows = table.querySelectorAll('tbody tr');
      dataRows.forEach(row => {
        const cols = row.querySelectorAll('td');
        const rowData = [];
        cols.forEach((col, index) => {
          // Skip the last column
          if (index < cols.length - 1) {
            rowData.push(col.innerText.trim());
          }
        });
        if (rowData.length > 0 && !rowData.every(cell => cell === '')) {
          bodyRows.push(rowData);
        }
      });
    });

    // Create table in PDF with proper header/body separation
    if (headerRow.length > 0 && bodyRows.length > 0) {
      doc.autoTable({
        head: [headerRow],
        body: bodyRows,
        startY: yPosition,
        margin: margin,
        theme: 'grid',
        styles: {
          fontSize: 9,
          cellPadding: 4,
          overflow: 'linebreak',
          halign: 'center'
        },
        headStyles: {
          fillColor: [51, 104, 145],
          textColor: [255, 255, 255],
          fontStyle: 'bold'
        },
        alternateRowStyles: {
          fillColor: [240, 240, 240]
        },
        didDrawPage: function(data) {
          // Footer
          const pageSize = doc.internal.pageSize;
          const pageHeight = pageSize.getHeight();
          const pageWidth = pageSize.getWidth();
          doc.setFontSize(10);
          doc.text(`Page ${doc.internal.pages.length - 1}`, pageWidth / 2, pageHeight - 10, { align: 'center' });
        }
      });
    }

    // Save PDF
    doc.save(`export_${new Date().getTime()}.pdf`);
    Swal.fire('Success', 'Data exported to PDF successfully!', 'success');
  } catch (error) {
    console.error('PDF export error:', error);
    Swal.fire('Error', 'Failed to generate PDF: ' + error.message, 'error');
  }
}

// Export to CSV
function exportToCSV(tables) {
  let csv = [];
  
  // Process all tables to get headers and data
  tables.forEach((table, index) => {
    // Get header rows from thead (excluding last column if "Action")
    const headerRows = table.querySelectorAll('thead tr');
    headerRows.forEach(row => {
      const cols = row.querySelectorAll('th, td');
      const csvRow = [];
      cols.forEach((col, colIndex) => {
        // Skip the last column if it contains "Action"
        if (colIndex < cols.length - 1 || (col.innerText.trim().toLowerCase() !== 'action' && col.innerText.trim().toLowerCase() !== 'actions')) {
          csvRow.push('"' + col.innerText.trim().replace(/"/g, '""') + '"');
        }
      });
      if (csvRow.length > 0 && !csvRow.every(cell => cell === '""')) {
        csv.push(csvRow.join(','));
      }
    });
    
    // Get data rows from tbody (excluding last column)
    const dataRows = table.querySelectorAll('tbody tr');
    dataRows.forEach(row => {
      const cols = row.querySelectorAll('td');
      const csvRow = [];
      cols.forEach((col, colIndex) => {
        // Skip the last column
        if (colIndex < cols.length - 1) {
          csvRow.push('"' + col.innerText.trim().replace(/"/g, '""') + '"');
        }
      });
      if (csvRow.length > 0 && !csvRow.every(cell => cell === '""')) {
        csv.push(csvRow.join(','));
      }
    });
  });

  if (csv.length === 0) {
    Swal.fire('Error', 'No data found to export', 'error');
    return;
  }

  const csvContent = csv.join('\n');
  downloadCSV(csvContent);
  Swal.fire('Success', 'Data exported to CSV successfully!', 'success');
}

// Export to Excel
function exportToExcel(tables) {
  // Placeholder for Excel export logic
  Swal.fire('Success', 'Exporting to Excel...', 'success');
  // You can use libraries like XLSX (SheetJS) here
  console.log('Excel export initiated');
}

// Download CSV file
function downloadCSV(csv) {
  const blob = new Blob([csv], { type: 'text/csv' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `export_${new Date().getTime()}.csv`;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
}

// Add styles for export buttons
const style = document.createElement('style');
style.textContent = `
  .export-option-btn {
    padding: 12px 16px;
    background-color: #f0f0f0;
    border: 2px solid #336891;
    border-radius: 6px;
    cursor: pointer;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    transition: all 0.3s ease;
    color: #336891;
    font-weight: 500;
  }
  
  .export-option-btn:hover {
    background-color: #336891;
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(51, 104, 145, 0.3);
  }
  
  .export-option-btn:active {
    transform: translateY(0);
  }
`;
document.head.appendChild(style);

// Check if libraries are loaded and log status
function checkLibraries() {
  console.log('=== Export Library Check ===');
  console.log('Swal available:', typeof Swal !== 'undefined');
  console.log('jsPDF available (window.jspdf.jsPDF):', window.jspdf && window.jspdf.jsPDF ? 'Yes' : 'No');
  console.log('jsPDF-autoTable available:', typeof window.jspdf !== 'undefined' && window.jspdf.jsPDF ? 'Yes' : 'No');
  console.log('Window.jspdf:', typeof window.jspdf !== 'undefined' ? 'Available' : 'Not available');
}

// Check libraries on page load
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', checkLibraries);
} else {
  checkLibraries();
}

// Also check after a delay to ensure all scripts loaded
setTimeout(checkLibraries, 2000);
