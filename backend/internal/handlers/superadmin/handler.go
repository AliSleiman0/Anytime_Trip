package superadmin

import (
	models "travel/backend/internal/models/superadmin"
	"travel/backend/internal/repository/superadmin"
	"bytes"
	"fmt"
	"html/template"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// SuperAdminHandler handles super admin-level requests
type SuperAdminHandler struct {
	configRepo *superadmin.SystemConfigRepository
	answerRepo *superadmin.PredefinedAnswerRepository
}

func NewSuperAdminHandler(configRepo *superadmin.SystemConfigRepository, answerRepo *superadmin.PredefinedAnswerRepository) *SuperAdminHandler {
	return &SuperAdminHandler{
		configRepo: configRepo,
		answerRepo: answerRepo,
	}
}

func (h *SuperAdminHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Super Admin Dashboard",
	})
}

func (h *SuperAdminHandler) ManageSystem(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Super Admin System Management",
	})
}

// Analytics handlers
func (h *SuperAdminHandler) GetAnalytics(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/analytics.html")
}

func (h *SuperAdminHandler) GetAnalyticsFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/analytics-frag.html")
}

// Accounting handlers
func (h *SuperAdminHandler) GetAccounting(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/accounting.html")
}

func (h *SuperAdminHandler) GetAccountingFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/accounting-frag.html")
}

// Manage Agents handlers
func (h *SuperAdminHandler) GetManageAgents(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/manage-agents.html")
}

func (h *SuperAdminHandler) GetManageAgentsFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/manage-agents-frag.html")
}

// Predefined Answers handlers
func (h *SuperAdminHandler) GetPredefinedAnswers(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/predefined-answers.html")
}

func (h *SuperAdminHandler) GetPredefinedAnswersFragment(c *fiber.Ctx) error {
	answers, err := h.answerRepo.FindAll(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error loading answers")
	}

	tmpl := `
<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 8px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 10px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #336891;
    border-radius: 10px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #2a5573;
  }
</style>

<div class="fixed top-0 left-64 right-0 h-screen overflow-auto custom-scrollbar">
  <main>
    <div hx-get="/admin/header?title=Predefined Answers&subtitle=Manage customer support responses" hx-trigger="load" hx-swap="innerHTML"></div>

    <div class="p-6">
      <div>
        <button hx-post="/superadmin/predefined-answers" hx-target="#answers-container" hx-swap="beforeend" class="bg-[#336891] w-126 h-12 text-white text-lg rounded-xl flex justify-between items-center p-2">
          <span class="inline-flex items-center justify-center w-8 h-8">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M18 12.998H13V17.998C13 18.2633 12.8946 18.5176 12.7071 18.7052C12.5196 18.8927 12.2652 18.998 12 18.998C11.7348 18.998 11.4804 18.8927 11.2929 18.7052C11.1054 18.5176 11 18.2633 11 17.998V12.998H6C5.73478 12.998 5.48043 12.8927 5.29289 12.7052C5.10536 12.5176 5 12.2633 5 11.998C5 11.7328 5.10536 11.4785 5.29289 11.2909C5.48043 11.1034 5.73478 10.998 6 10.998H11V5.99805C11 5.73283 11.1054 5.47848 11.2929 5.29094C11.4804 5.1034 11.7348 4.99805 12 4.99805C12.2652 4.99805 12.5196 5.1034 12.7071 5.29094C12.8946 5.47848 13 5.73283 13 5.99805V10.998H18C18.2652 10.998 18.5196 11.1034 18.7071 11.2909C18.8946 11.4785 19 11.7328 19 11.998C19 12.2633 18.8946 12.5176 18.7071 12.7052C18.5196 12.8927 18.2652 12.998 18 12.998Z" fill="white"/>
            </svg>
          </span>
          <span>Add New Predefined Answer</span>
        </button>
      </div>

      <div id="answers-container">
        {{range .Answers}}
        <div class="answer-section" id="answer-{{.ID}}">
          <div class="flex items-center py-5">
            <div class="mr-4 bg-[#F8F8F8] p-2 rounded-lg shadow-md w-1/6 h-20 flex flex-col justify-center items-center text-[#6B6B6B] cursor-pointer hover:bg-gray-200 transition-colors" onclick="openShortcutBindingPopup({{.ID}}, '{{.Shortcut}}')" data-answer-id="{{.ID}}">
              <p id="shortcut-display-{{.ID}}">{{.Shortcut}}</p>
            </div>
            <div class="bg-[#F8F8F8] p-2 px-5 rounded-lg shadow-md w-5/6 h-20 flex items-center justify-between">
              <p class="answer-text">{{.Answer}}</p>

              <div class="flex gap-2.5">
                <button onclick="toggleEdit({{.ID}})" class="bg-[#336891] w-5 h-5 flex items-center justify-center rounded-md">
                  <svg width="11" height="11" viewBox="0 0 11 11" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M9.2155 0.341584C8.99672 0.122868 8.70002 0 8.39067 0C8.08131 0 7.78462 0.122868 7.56583 0.341584L6.60625 1.30175L9.69967 4.39517L10.6593 3.43617C10.7676 3.32783 10.8536 3.19919 10.9123 3.05762C10.9709 2.91604 11.0011 2.76429 11.0011 2.61104C11.0011 2.45779 10.9709 2.30605 10.9123 2.16447C10.8536 2.02289 10.7676 1.89426 10.6593 1.78592L9.2155 0.341584ZM8.87483 5.22L5.78142 2.12658L0.627667 7.28033L0 11.002L3.72167 10.3738L8.87483 5.22Z" fill="white"/>
                  </svg>
                </button>
                <button hx-delete="/superadmin/predefined-answers/{{.ID}}" hx-target="#answer-{{.ID}}" hx-swap="outerHTML" hx-confirm="Are you sure you want to delete this answer?" class="bg-[#D24124] w-5 h-5 flex items-center justify-center rounded-md">
                  <svg width="9" height="11" viewBox="0 0 9 11" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M0.583333 9.33333C0.583333 9.975 1.10833 10.5 1.75 10.5H6.41667C7.05833 10.5 7.58333 9.975 7.58333 9.33333V3.5C7.58333 2.85833 7.05833 2.33333 6.41667 2.33333H1.75C1.10833 2.33333 0.583333 2.85833 0.583333 3.5V9.33333ZM7.58333 0.583333H6.125L5.71083 0.169167C5.60583 0.0641666 5.45417 0 5.3025 0H2.86417C2.7125 0 2.56083 0.0641666 2.45583 0.169167L2.04167 0.583333H0.583333C0.2625 0.583333 0 0.845833 0 1.16667C0 1.4875 0.2625 1.75 0.583333 1.75H7.58333C7.90417 1.75 8.16667 1.4875 8.16667 1.16667C8.16667 0.845833 7.90417 0.583333 7.58333 0.583333Z" fill="white"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div>
            <input type="text" value="{{.Answer}}" id="edit-{{.ID}}" class="w-full p-3 bg-[#F8F8F8] h-20 hidden rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500" onkeypress="if(event.key === 'Enter') saveEdit({{.ID}})" />
          </div>
        </div>
        {{end}}
      </div>
    </div>
  </main>
</div>

<script>
let capturedKeys = [];
let isCapturingShortcut = false;
let currentAnswerId = null;
let pressedKeys = new Set();

function openShortcutBindingPopup(answerId, currentShortcut) {
  currentAnswerId = answerId;
  capturedKeys = [];
  pressedKeys.clear();
  isCapturingShortcut = true;
  
  Swal.fire({
    title: 'Enter a Key Bind',
    html: ` + "`" + `
      <div class="text-center">
        <p class="text-gray-600 mb-4">Press two keys together to set as shortcut</p>
        <div id="captured-keys" class="text-2xl font-bold text-blue-600 min-h-[40px]">
          Waiting for input...
        </div>
        <p class="text-sm text-gray-500 mt-2">Current: ${currentShortcut || 'None'}</p>
      </div>
    ` + "`" + `,
    showCancelButton: true,
    showConfirmButton: false,
    cancelButtonText: 'Cancel',
    cancelButtonColor: '#336891',
    customClass: {
      cancelButton: 'swal-cancel-btn',
      popup: 'swal-popup'
    },
    didOpen: () => {
      document.addEventListener('keydown', handleShortcutCapture, true);
      document.addEventListener('keyup', handleKeyRelease, true);
    },
    willClose: () => {
      document.removeEventListener('keydown', handleShortcutCapture, true);
      document.removeEventListener('keyup', handleKeyRelease, true);
      isCapturingShortcut = false;
      capturedKeys = [];
      pressedKeys.clear();
    }
  }).catch(() => {});
}

function handleKeyRelease(event) {
  if (!isCapturingShortcut) return;
  pressedKeys.clear();
}

function handleShortcutCapture(event) {
  if (!isCapturingShortcut) return;
  
  event.preventDefault();
  event.stopPropagation();
  
  if (event.repeat) return;
  
  let shortcut = [];
  if (event.ctrlKey) shortcut.push('Ctrl');
  if (event.altKey) shortcut.push('Alt');
  if (event.shiftKey) shortcut.push('Shift');
  if (event.metaKey) shortcut.push('Meta');
  
  if (['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) {
    return;
  }
  
  let mainKey = '';
  if (event.key.length === 1) {
    mainKey = event.key.toUpperCase();
  } else if (event.code.startsWith('Digit')) {
    mainKey = event.code.replace('Digit', '');
  } else if (event.code.startsWith('Key')) {
    mainKey = event.code.replace('Key', '');
  } else if (event.code.startsWith('Arrow')) {
    mainKey = event.code.replace('Arrow', '');
  } else if (['Enter', 'Space', 'Tab', 'Escape', 'Backspace'].includes(event.key)) {
    mainKey = event.key;
  } else {
    mainKey = event.key;
  }
  
  shortcut.push(mainKey);
  
  if (shortcut.length < 2) {
    document.getElementById('captured-keys').textContent = 'Please press a modifier key (Ctrl, Alt, Shift) + another key';
    return;
  }
  
  const shortcutString = shortcut.join(' + ');
  document.getElementById('captured-keys').textContent = shortcutString;
  
  setTimeout(() => {
    if (isCapturingShortcut) {
      saveShortcut(currentAnswerId, shortcutString);
    }
  }, 300);
}

function saveShortcut(answerId, shortcut) {
  isCapturingShortcut = false;
  
  fetch(` + "`/superadmin/predefined-answers/${answerId}/shortcut`" + `, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ shortcut: shortcut })
  })
  .then(response => {
    if (!response.ok) throw new Error('Server responded with error');
    return response.json();
  })
  .then(data => {
    if (data.success) {
      const displayElement = document.getElementById(` + "`shortcut-display-${answerId}`" + `);
      if (displayElement) {
        displayElement.textContent = shortcut;
      }
      
      Swal.close();
      Swal.fire({
        title: 'Success!',
        text: ` + "`Shortcut set to: ${shortcut}`" + `,
        icon: 'success',
        timer: 2000,
        showConfirmButton: false
      });
    } else {
      Swal.fire({
        title: 'Error',
        text: 'Failed to save shortcut',
        icon: 'error'
      });
    }
  })
  .catch(error => {
    console.error('Error saving shortcut:', error);
    Swal.fire({
      title: 'Error',
      text: 'Failed to save shortcut',
      icon: 'error'
    });
  });
}

function toggleEdit(id) {
  const editInput = document.getElementById('edit-' + id);
  const answerText = document.querySelector('#answer-' + id + ' .answer-text');
  
  if (editInput.classList.contains('hidden')) {
    editInput.classList.remove('hidden');
    editInput.focus();
  } else {
    saveEdit(id);
  }
}

function saveEdit(id) {
  const editInput = document.getElementById('edit-' + id);
  const newValue = editInput.value;
  
  fetch('/superadmin/predefined-answers/' + id, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ answer: newValue })
  }).then(response => {
    if (response.ok) {
      const answerText = document.querySelector('#answer-' + id + ' .answer-text');
      answerText.textContent = newValue;
      editInput.classList.add('hidden');
    }
  });
}
</script>
`

	t, err := template.New("predefined-answers").Parse(tmpl)
	if err != nil {
		return c.Status(500).SendString("Template error")
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, fiber.Map{"Answers": answers})
	if err != nil {
		return c.Status(500).SendString("Render error")
	}

	return c.Type("html").Send(buf.Bytes())
}

// CreatePredefinedAnswer creates a new predefined answer
func (h *SuperAdminHandler) CreatePredefinedAnswer(c *fiber.Ctx) error {
	// Get the current count to generate the next shortcut number
	answers, err := h.answerRepo.FindAll(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error loading answers: " + err.Error())
	}
	nextNum := len(answers) + 1

	newAnswer := &models.PredefinedAnswer{
		ID:       nextNum,
		Shortcut: fmt.Sprintf("Ctrl + %d", nextNum),
		Answer:   "Predefined Answer " + strconv.Itoa(nextNum),
	}

	err = h.answerRepo.Create(c.Context(), newAnswer)
	if err != nil {
		return c.Status(500).SendString("Error creating answer: " + err.Error())
	}

	// Return the new section HTML
	tmpl := `
<div class="answer-section" id="answer-{{.ID}}">
  <div class="flex items-center py-5">
    <div class="mr-4 bg-[#F8F8F8] p-2 rounded-lg shadow-md w-1/6 h-20 flex flex-col justify-center items-center text-[#6B6B6B] cursor-pointer hover:bg-gray-200 transition-colors" onclick="openShortcutBindingPopup({{.ID}}, '{{.Shortcut}}')" data-answer-id="{{.ID}}">
      <p id="shortcut-display-{{.ID}}">{{.Shortcut}}</p>
    </div>
    <div class="bg-[#F8F8F8] p-2 px-5 rounded-lg shadow-md w-5/6 h-20 flex items-center justify-between">
      <p class="answer-text">{{.Answer}}</p>

      <div class="flex gap-2.5">
        <button onclick="toggleEdit({{.ID}})" class="bg-[#336891] w-5 h-5 flex items-center justify-center rounded-md">
          <svg width="11" height="11" viewBox="0 0 11 11" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M9.2155 0.341584C8.99672 0.122868 8.70002 0 8.39067 0C8.08131 0 7.78462 0.122868 7.56583 0.341584L6.60625 1.30175L9.69967 4.39517L10.6593 3.43617C10.7676 3.32783 10.8536 3.19919 10.9123 3.05762C10.9709 2.91604 11.0011 2.76429 11.0011 2.61104C11.0011 2.45779 10.9709 2.30605 10.9123 2.16447C10.8536 2.02289 10.7676 1.89426 10.6593 1.78592L9.2155 0.341584ZM8.87483 5.22L5.78142 2.12658L0.627667 7.28033L0 11.002L3.72167 10.3738L8.87483 5.22Z" fill="white"/>
          </svg>
        </button>
        <button hx-delete="/superadmin/predefined-answers/{{.ID}}" hx-target="#answer-{{.ID}}" hx-swap="outerHTML" hx-confirm="Are you sure you want to delete this answer?" class="bg-[#D24124] w-5 h-5 flex items-center justify-center rounded-md">
          <svg width="9" height="11" viewBox="0 0 9 11" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M0.583333 9.33333C0.583333 9.975 1.10833 10.5 1.75 10.5H6.41667C7.05833 10.5 7.58333 9.975 7.58333 9.33333V3.5C7.58333 2.85833 7.05833 2.33333 6.41667 2.33333H1.75C1.10833 2.33333 0.583333 2.85833 0.583333 3.5V9.33333ZM7.58333 0.583333H6.125L5.71083 0.169167C5.60583 0.0641666 5.45417 0 5.3025 0H2.86417C2.7125 0 2.56083 0.0641666 2.45583 0.169167L2.04167 0.583333H0.583333C0.2625 0.583333 0 0.845833 0 1.16667C0 1.4875 0.2625 1.75 0.583333 1.75H7.58333C7.90417 1.75 8.16667 1.4875 8.16667 1.16667C8.16667 0.845833 7.90417 0.583333 7.58333 0.583333Z" fill="white"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
  <div>
    <input type="text" value="{{.Answer}}" id="edit-{{.ID}}" class="w-full p-3 bg-[#F8F8F8] h-20 hidden rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500" onkeypress="if(event.key === 'Enter') saveEdit({{.ID}})" />
  </div>
</div>
`

	t, err := template.New("answer").Parse(tmpl)
	if err != nil {
		return c.Status(500).SendString("Template error: " + err.Error())
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, newAnswer)
	if err != nil {
		return c.Status(500).SendString("Render error: " + err.Error())
	}

	return c.Type("html").Send(buf.Bytes())
}

// UpdatePredefinedAnswer updates a predefined answer
func (h *SuperAdminHandler) UpdatePredefinedAnswer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	var req struct {
		Answer string `json:"answer"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	answer, err := h.answerRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(404).SendString("Answer not found")
	}

	answer.Answer = req.Answer
	err = h.answerRepo.Update(c.Context(), id, answer)
	if err != nil {
		return c.Status(500).SendString("Error updating answer")
	}

	return c.SendStatus(200)
}

// UpdatePredefinedAnswerShortcut updates a predefined answer shortcut
func (h *SuperAdminHandler) UpdatePredefinedAnswerShortcut(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	var req struct {
		Shortcut string `json:"shortcut"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	answer, err := h.answerRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(404).SendString("Answer not found")
	}

	answer.Shortcut = req.Shortcut
	err = h.answerRepo.Update(c.Context(), id, answer)
	if err != nil {
		return c.Status(500).SendString("Error updating shortcut")
	}

	return c.JSON(fiber.Map{"success": true})
}

// DeletePredefinedAnswer deletes a predefined answer
func (h *SuperAdminHandler) DeletePredefinedAnswer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	err = h.answerRepo.Delete(c.Context(), id)
	if err != nil {
		return c.Status(500).SendString("Error deleting answer")
	}

	// Return empty response to remove the element
	return c.SendString("")
}
