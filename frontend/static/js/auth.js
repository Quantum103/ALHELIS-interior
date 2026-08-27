/* =========================================
   1. Dynamic Background & Visual FX
   ========================================= */
const cg = document.getElementById('cg');
let mx = innerWidth / 2, my = innerHeight / 2, cx = mx, cy = my;

addEventListener('mousemove', e => { mx = e.clientX; my = e.clientY; });
addEventListener('mouseleave', () => { if (cg) cg.style.opacity = '0'; });
addEventListener('mouseenter', () => { if (cg) cg.style.opacity = '1'; });

(function loop() {
  cx += (mx - cx) * 0.08;
  cy += (my - cy) * 0.08;
  if (cg) {
    cg.style.left = cx + 'px';
    cg.style.top = cy + 'px';
  }
  requestAnimationFrame(loop);
})();

const mm = matchMedia('(min-width:1024px)');
function initParticles() {
  if (!mm.matches) return;
  const c = document.getElementById('ptc');
  if (!c || c.childElementCount) return;
  for (let i = 0; i < 22; i++) {
    const p = document.createElement('span');
    p.className = 'particle';
    p.style.left = Math.random() * 100 + '%';
    p.style.bottom = '-10px';
    p.style.animationDuration = (10 + Math.random() * 18) + 's';
    p.style.animationDelay = (Math.random() * 15) + 's';
    p.style.opacity = Math.random() * 0.6 + 0.2;
    c.appendChild(p);
  }
}
mm.addEventListener('change', initParticles);
initParticles();

/* =========================================
   2. Form Switcher & Password Mechanics
   ========================================= */
const sw = document.getElementById('sw');
const btns = sw ? [...sw.querySelectorAll('button')] : [];
const loginF = document.getElementById('loginForm');
const regF = document.getElementById('regForm');
const fTitle = document.getElementById('fTitle');
const fSub = document.getElementById('fSub');

btns.forEach(b => b.addEventListener('click', () => {
  const m = b.dataset.m;
  btns.forEach(x => x.classList.toggle('on', x === b));
  sw.classList.toggle('right', m === 'register');
  if (m === 'register') {
    loginF.classList.add('hidden');
    regF.classList.remove('hidden');
    fTitle.innerHTML = 'Создать <em>аккаунт</em>';
    fSub.textContent = 'Начните работу со студией';
  } else {
    regF.classList.add('hidden');
    loginF.classList.remove('hidden');
    fTitle.innerHTML = 'Добро <em>пожаловать</em>';
    fSub.textContent = 'Войдите, чтобы продолжить';
  }
}));

document.querySelectorAll('.pwd-toggle').forEach(btn => {
  btn.addEventListener('click', () => {
    const inp = document.getElementById(btn.dataset.for);
    if (!inp) return;
    const show = inp.type === 'password';
    inp.type = show ? 'text' : 'password';
    btn.textContent = show ? 'Скрыть' : 'Показать';
  });
});

const rPwd = document.getElementById('rPwd');
const str = document.getElementById('str');
const hint = document.getElementById('hint');

function scorePassword(v) {
  let s = 0;
  if (v.length >= 6) s++;
  if (v.length >= 10) s++;
  if (/[A-Z]/.test(v) && /[a-z]/.test(v)) s++;
  if (/\d/.test(v) && /[^A-Za-z0-9]/.test(v)) s++;
  return s;
}

if (rPwd && str && hint) {
  rPwd.addEventListener('input', () => {
    const v = rPwd.value;
    const s = v ? scorePassword(v) : 0;
    str.className = 'strength' + (s ? ' w' + s : '');
    const msgs = ['Введите пароль', 'Слабый', 'Средний', 'Хороший', 'Отличный'];
    hint.textContent = msgs[s] || msgs[0];
    hint.style.color = ['var(--muted)', '#c74a4a', 'var(--bronze)', 'var(--gold)', '#8fa876'][s];
  });
}

/* =========================================
   3. Helpers & Validation
   ========================================= */
function validEmail(v) { 
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v); 
}

function setErr(inp, on) { 
  if (inp) inp.closest('.field').classList.toggle('err', on); 
}

function pulseCard(card) {
  if (!card) return;
  card.animate([
    { boxShadow: '0 0 0 0 rgba(214,184,138,.55)' },
    { boxShadow: '0 0 0 30px rgba(214,184,138,0)' }
  ], { duration: 800, easing: 'cubic-bezier(.2,.8,.2,1)' });
}

function showSuccessCard(message, redirectUrl) {
  const card = document.getElementById('card');
  if (!card) return;

  card.innerHTML = `
    <div class="success">
      <div class="ic">✓</div>
      <h3>${message}</h3>
      <p>Перенаправляем в кабинет…</p>
      <a class="submit" style="margin-top:26px;display:inline-block;max-width:260px;text-decoration:none" href="${redirectUrl || '/dashboard'}">Перейти →</a>
    </div>`;

  if (redirectUrl) {
    setTimeout(() => { window.location.href = redirectUrl; }, 1200);
  }
}

/* =========================================
   4. API Handlers & Form Submit
   ========================================= */
function setupFormSubmit(form, btn, validateAndGetPayload, apiPath, successMsg) {
  if (!form || !btn) return;

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    // 1. Проверка и сбор данных из формы
    const payload = validateAndGetPayload();
    if (!payload) return; // Ошибка валидации — останавливаем отправку

    // 2. Индикация загрузки
    btn.classList.add('loading');
    btn.disabled = true;
    pulseCard(document.getElementById('card'));

    try {
      // 3. Отправка запроса через Nginx Gateway
      const response = await fetch(apiPath, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Ошибка сервера при обработке запроса');
      }

      // 4. Успешный Вход
      if (apiPath.includes('/login')) {
        if (data.access_token) localStorage.setItem('access_token', data.access_token);
        if (data.token) localStorage.setItem('token', data.token); // Для обратной совместимости
        if (data.user) localStorage.setItem('user', JSON.stringify(data.user));

        showSuccessCard(successMsg, '/dashboard');
      } 
      // 5. Успешная Регистрация
      else if (apiPath.includes('/register')) {
        showSuccessCard(successMsg + '<br><small style="font-size:14px;font-weight:400;opacity:0.8">Вы можете войти под своим логином</small>', null);
        setTimeout(() => {
          location.reload(); // Перезагрузка страницы для входа
        }, 2000);
      }

    } catch (err) {
      btn.classList.remove('loading');
      btn.disabled = false;
      alert(err.message);
    }
  });
}

// Сброс красной подсветки при вводе
document.querySelectorAll('.field input').forEach(inp => {
  inp.addEventListener('input', () => setErr(inp, false));
});

// Настройка формы ВХОДА
setupFormSubmit(
  loginF, 
  document.getElementById('lSub'), 
  () => {
    const u = document.getElementById('lUsername');
    const p = document.getElementById('lPwd');
    
    const uOk = u.value.trim().length >= 2;
    const pOk = p.value.length >= 6;

    setErr(u, !uOk); 
    setErr(p, !pOk);

    if (!uOk) { u.focus(); return null; }
    if (!pOk) { p.focus(); return null; }

    return { 
      username: u.value.trim(), 
      password: p.value 
    };
  }, 
  '/api/auth/login', 
  'Вы <em>в системе</em>'
);

// Настройка формы РЕГИСТРАЦИИ
setupFormSubmit(
  regF, 
  document.getElementById('rSub'), 
  () => {
    const n = document.getElementById('rName');
    const e = document.getElementById('rEmail');
    const p = document.getElementById('rPwd');
    const p2 = document.getElementById('rPwd2');
    const ag = document.getElementById('agree');

    const nOk = n.value.trim().length >= 2;
    const eOk = validEmail(e.value.trim());
    const pOk = p.value.length >= 6;
    const p2Ok = p2.value && p2.value === p.value;

    setErr(n, !nOk); 
    setErr(e, !eOk); 
    setErr(p, !pOk); 
    setErr(p2, !p2Ok);

    if (!nOk) { n.focus(); return null; }
    if (!eOk) { e.focus(); return null; }
    if (!pOk) { p.focus(); return null; }
    if (!p2Ok) { p2.focus(); return null; }

    if (!ag.checked) {
      ag.parentElement.animate([
        { transform: 'translateX(0)' }, 
        { transform: 'translateX(-6px)' }, 
        { transform: 'translateX(6px)' }, 
        { transform: 'translateX(0)' }
      ], { duration: 300 });
      return null;
    }

    return { 
      name: n.value.trim(), 
      email: e.value.trim(), 
      password: p.value 
    };
  }, 
  '/api/auth/register', 
  'Аккаунт <em>создан</em>'
);