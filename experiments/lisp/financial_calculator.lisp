(defstruct financial-params
  (num-years 0 :type integer)
  (au-hours 0.0 :type float)
  (initial-tsn 0.0 :type float)
  (rate-escalation 0.0 :type float)
  (aic 0.0 :type float)
  (hsi-tsn 0.0 :type float)
  (overhaul-tsn 0.0 :type float)
  (hsi-cost 0.0 :type float)
  (overhaul-cost 0.0 :type float))

(defun validate-params (params)
  (flet ((check (condition error-message)
           (when condition (error error-message))))
    (check (<= (financial-params-num-years params) 0) "NumYears must be positive")
    (check (<= (financial-params-au-hours params) 0.0) "AuHours must be positive")
    (check (< (financial-params-initial-tsn params) 0.0) "InitialTSN cannot be negative")
    (check (< (financial-params-rate-escalation params) 0.0) "RateEscalation cannot be negative")
    (check (or (< (financial-params-aic params) 0.0) (> (financial-params-aic params) 100.0))
           "AIC must be between 0 and 100")
    (check (<= (financial-params-hsi-tsn params) 0.0) "HSITSN must be positive")
    (check (<= (financial-params-overhaul-tsn params) 0.0) "OverhaulTSN must be positive")
    (check (< (financial-params-hsi-cost params) 0.0) "HSICost cannot be negative")
    (check (< (financial-params-overhaul-cost params) 0.0) "OverhaulCost cannot be negative")))

(defun calculate-financials (rate params)
  (let ((cumulative-profit 0.0))
    (dotimes (year (financial-params-num-years params) cumulative-profit)
      (let* ((tsn (+ (financial-params-initial-tsn params)
                     (* (financial-params-au-hours params) (1+ year))))
             (escalated-rate (* rate (expt (1+ (/ (financial-params-rate-escalation params) 100.0)) year)))
             (engine-revenue (* (financial-params-au-hours params) escalated-rate))
             (aic-revenue (* engine-revenue (/ (financial-params-aic params) 100.0)))
             (total-revenue (+ engine-revenue aic-revenue))
             (hsi (and (>= tsn (financial-params-hsi-tsn params))
                       (or (= year 0)
                           (< (- tsn (financial-params-au-hours params))
                              (financial-params-hsi-tsn params)))))
             (overhaul (and (>= tsn (financial-params-overhaul-tsn params))
                            (or (= year 0)
                                (< (- tsn (financial-params-au-hours params))
                                   (financial-params-overhaul-tsn params)))))
             (hsi-cost (if hsi (financial-params-hsi-cost params) 0.0))
             (overhaul-cost (if overhaul (financial-params-overhaul-cost params) 0.0))
             (total-cost (+ hsi-cost overhaul-cost))
             (total-profit (- total-revenue total-cost)))
        (incf cumulative-profit total-profit)))))
;; ... [Previous code remains the same up to the newton-raphson function] ...

(defun newton-raphson (f df x0 xtol max-iter)
  (loop for i from 0 below max-iter
        for fx = (funcall f x0)
        for dfx = (funcall df x0)
        do (cond ((< (abs fx) xtol) (return-from newton-raphson (values x0 (1+ i))))
                 ((zerop dfx) (error "Derivative is zero, can't proceed with Newton-Raphson"))
                 (t (setf x0 (- x0 (/ fx dfx)))))
        finally (error "Newton-Raphson method did not converge within ~D iterations" max-iter)))

(defun goal-seek (target-profit params initial-guess)
  (labels ((objective (rate)
             (- (calculate-financials rate params) target-profit))
           (derivative (rate)
             (let ((epsilon 1e-6))
               (/ (- (objective (+ rate epsilon))
                     (objective rate))
                  epsilon))))
    (newton-raphson #'objective #'derivative initial-guess 1e-8 100)))

(defun main ()
  (let ((params (make-financial-params :num-years 10
                                       :au-hours 450.0
                                       :initial-tsn 100.0
                                       :rate-escalation 5.0
                                       :aic 10.0
                                       :hsi-tsn 1000.0
                                       :overhaul-tsn 3000.0
                                       :hsi-cost 50000.0
                                       :overhaul-cost 100000.0))
        (initial-rate 100.0)
        (target-profit 3000000.0))
    
    (handler-case (validate-params params)
      (error (e) (format t "Invalid parameters: ~A~%" e) (return-from main)))
    
    (let ((start-time (get-internal-real-time)))
      (handler-case
          (let ((initial-cumulative-profit (calculate-financials initial-rate params)))
            (format t "Initial Warranty Rate: ~,2F~%" initial-rate)
            (format t "Initial Cumulative Profit: ~,2F~%" initial-cumulative-profit)
            
            (multiple-value-bind (optimal-rate iterations)
                (goal-seek target-profit params initial-rate)
              (format t "~%Optimal Warranty Rate to achieve ~,2F profit: ~,7F~%"
                      target-profit optimal-rate)
              (format t "Number of iterations: ~D~%" iterations)
              
              (let ((final-cumulative-profit (calculate-financials optimal-rate params)))
                (format t "~%Final Cumulative Profit: ~,2F~%" final-cumulative-profit))))
        (error (e) (format t "Error: ~A~%" e)))
      
      (let ((end-time (get-internal-real-time)))
        (format t "~%Execution time: ~,3F seconds~%"
                (/ (- end-time start-time) internal-time-units-per-second))))))

;; Run the main function
(main)
